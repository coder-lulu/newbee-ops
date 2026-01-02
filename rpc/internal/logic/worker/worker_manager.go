package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxymetrics"
	"github.com/zeromicro/go-zero/core/logx"
)

// WorkerManagerConfig Worker管理器配置
type WorkerManagerConfig struct {
	TTL                  time.Duration // Worker TTL（默认3分钟）
	HeartbeatTimeout     time.Duration // 心跳超时（默认2分钟）
	UpdateBatchSize      int           // 批量更新大小（默认100）
	UpdateInterval       time.Duration // 更新间隔（默认30秒）
	MetricsRetention     time.Duration // 指标保留时间（默认7天）
	HealthCheckInterval  time.Duration // 健康检查间隔（默认30秒）
	HealthCheckTimeout   time.Duration // 健康检查超时（默认5秒）
	DefaultStrategy      string        // 默认选择策略（默认least_connections）
}

// DefaultWorkerManagerConfig 默认配置
func DefaultWorkerManagerConfig() *WorkerManagerConfig {
	return &WorkerManagerConfig{
		TTL:                  3 * time.Minute,
		HeartbeatTimeout:     2 * time.Minute,
		UpdateBatchSize:      100,
		UpdateInterval:       30 * time.Second,
		MetricsRetention:     7 * 24 * time.Hour,
		HealthCheckInterval:  30 * time.Second,
		HealthCheckTimeout:   5 * time.Second,
		DefaultStrategy:      "least_connections",
	}
}

// WorkerMetricsSnapshot 上一次指标快照（用于计算增量）
type WorkerMetricsSnapshot struct {
	Timestamp     time.Time
	NetworkIn     uint64
	NetworkOut    uint64
	TotalRequests uint64
	SuccessCount  uint64
	FailureCount  uint64
}

// WorkerManager Worker管理器 - 双层存储架构核心
type WorkerManager struct {
	config   *WorkerManagerConfig
	registry *WorkerRegistry
	db       *ent.Client
	logger   logx.Logger

	// 异步更新队列
	updateQueue chan *WorkerInfo
	updateWg    sync.WaitGroup

	// 后台任务控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// 选择策略
	strategies map[string]SelectionStrategy

	// 健康检查器
	healthChecker *HealthChecker

	// 上一次指标快照（用于计算增量）
	mu                sync.RWMutex
	previousMetrics   map[string]*WorkerMetricsSnapshot
	healthCheckLatency map[string]float64 // workerID -> 最近健康检查延迟(ms)
}

// NewWorkerManager 创建Worker管理器
func NewWorkerManager(db *ent.Client, config *WorkerManagerConfig) *WorkerManager {
	if config == nil {
		config = DefaultWorkerManagerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger := logx.WithContext(ctx)

	// 初始化健康检查器
	healthChecker := NewHealthChecker(config.HealthCheckTimeout, logger)

	manager := &WorkerManager{
		config:             config,
		registry:           NewWorkerRegistry(config.TTL, logger),
		db:                 db,
		logger:             logger,
		updateQueue:        make(chan *WorkerInfo, config.UpdateBatchSize*10),
		ctx:                ctx,
		cancel:             cancel,
		strategies:         make(map[string]SelectionStrategy),
		healthChecker:      healthChecker,
		previousMetrics:    make(map[string]*WorkerMetricsSnapshot),
		healthCheckLatency: make(map[string]float64),
	}

	// 注册所有选择策略
	manager.registerStrategies()

	return manager
}

// registerStrategies 注册所有选择策略
func (m *WorkerManager) registerStrategies() {
	m.strategies["round_robin"] = NewRoundRobinStrategy()
	m.strategies["least_connections"] = NewLeastConnectionsStrategy()
	m.strategies["weighted_round_robin"] = NewWeightedRoundRobinStrategy()
	m.strategies["consistent_hash"] = NewConsistentHashStrategy()
	m.strategies["geo_nearest"] = NewGeoNearestStrategy()
	m.strategies["priority"] = NewPriorityStrategy()
	m.strategies["random"] = NewRandomStrategy()

	m.logger.Infow("Registered selection strategies",
		logx.Field("strategies", []string{
			"round_robin",
			"least_connections",
			"weighted_round_robin",
			"consistent_hash",
			"geo_nearest",
			"priority",
			"random",
		}))
}

// Start 启动Worker管理器
func (m *WorkerManager) Start() error {
	m.logger.Info("Starting WorkerManager...")

	// 1. 从数据库恢复Worker数据
	if err := m.recoverFromDatabase(); err != nil {
		return fmt.Errorf("failed to recover workers from database: %w", err)
	}

	// 2. 启动异步更新任务
	m.wg.Add(1)
	go m.asyncUpdateWorker()

	// 3. 启动TTL清理任务
	m.wg.Add(1)
	go m.ttlCleanup()

	// 4. 启动指标收集任务
	m.wg.Add(1)
	go m.asyncCollectMetrics()

	// 5. 启动健康检查任务
	m.wg.Add(1)
	go m.healthCheckLoop()

	// 6. 启动旧指标清理任务
	m.wg.Add(1)
	go m.metricsCleanup()

	m.logger.Info("✅ WorkerManager started successfully")
	return nil
}

// Stop 停止Worker管理器
func (m *WorkerManager) Stop() {
	m.logger.Info("Stopping WorkerManager...")

	// 取消所有后台任务
	m.cancel()

	// 等待所有后台任务完成
	m.wg.Wait()

	// 关闭更新队列
	close(m.updateQueue)
	m.updateWg.Wait()

	m.logger.Info("✅ WorkerManager stopped")
}

// Register 注册Worker
func (m *WorkerManager) Register(info *WorkerInfo) error {
	// 注册到内存
	if err := m.registry.Register(info); err != nil {
		return err
	}

	// 异步推送到更新队列
	select {
	case m.updateQueue <- info:
		// 成功入队
	default:
		// 队列满，记录警告但不阻塞
		m.logger.Errorw("Update queue full, worker update may be delayed",
			logx.Field("worker_id", info.WorkerID),
			logx.Field("queue_size", len(m.updateQueue)))
	}

	return nil
}

// Heartbeat 处理Worker心跳
func (m *WorkerManager) Heartbeat(workerID string, tenantID uint64, status string, load *WorkerLoad) error {
	// 更新内存
	if err := m.registry.Heartbeat(workerID, tenantID, status, load); err != nil {
		return err
	}

	// 获取完整信息用于数据库更新
	workerInfo, exists := m.registry.Get(workerID)
	if !exists {
		return fmt.Errorf("worker not found after heartbeat: %s", workerID)
	}

	// 异步推送到更新队列
	select {
	case m.updateQueue <- workerInfo:
		// 成功入队
	default:
		// 队列满，跳过此次更新（心跳频繁，丢失一次影响不大）
	}

	return nil
}

// SelectWorker 选择Worker
func (m *WorkerManager) SelectWorker(ctx context.Context, tenantID uint64, strategyName string, opts *SelectOptions) (*WorkerInfo, error) {
	// 获取租户的Worker列表
	workers := m.registry.List(tenantID, false)

	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers for tenant %d", tenantID)
	}

	// 如果未指定策略，使用默认策略
	if strategyName == "" {
		strategyName = m.config.DefaultStrategy
	}

	// 获取策略
	strategy, exists := m.strategies[strategyName]
	if !exists {
		m.logger.Errorw("Unknown selection strategy, using default",
			logx.Field("requested_strategy", strategyName),
			logx.Field("default_strategy", m.config.DefaultStrategy))
		strategy = m.strategies[m.config.DefaultStrategy]
	}

	// 执行选择
	selected, err := strategy.Select(ctx, workers, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to select worker: %w", err)
	}

	m.logger.Infow("Worker selected",
		logx.Field("worker_id", selected.WorkerID),
		logx.Field("strategy", strategyName),
		logx.Field("tenant_id", tenantID),
		logx.Field("active_sessions", selected.Load.ActiveSessions),
		logx.Field("status", selected.Status))

	return selected, nil
}

// GetWorker 获取单个Worker信息
func (m *WorkerManager) GetWorker(workerID string) (*WorkerInfo, error) {
	info, exists := m.registry.Get(workerID)
	if !exists {
		return nil, fmt.Errorf("worker not found: %s", workerID)
	}
	return info, nil
}

// ListWorkers 获取Worker列表
func (m *WorkerManager) ListWorkers(tenantID uint64) []*WorkerInfo {
	return m.registry.List(tenantID, false)
}

// GetWorkerCount 获取Worker统计
func (m *WorkerManager) GetWorkerCount(tenantID uint64) map[string]int {
	return m.registry.Count(tenantID)
}

// recoverFromDatabase 从数据库恢复Worker数据
func (m *WorkerManager) recoverFromDatabase() error {
	m.logger.Info("Recovering workers from database...")

	ctx := context.Background()
	systemCtx := hooks.NewSystemContext(ctx)

	// 查询所有在线和降级状态的Worker
	workers, err := m.db.Proxy.Query().
		Where(
			proxy.Or(
				proxy.WorkerStatusEQ(proxy.WorkerStatusOnline),
				proxy.WorkerStatusEQ(proxy.WorkerStatusDegraded),
			),
		).
		All(systemCtx)

	if err != nil {
		return fmt.Errorf("failed to query workers: %w", err)
	}

	// 恢复到内存
	recovered := 0
	for _, w := range workers {
		info := m.convertEntToWorkerInfo(w)
		if err := m.registry.Register(info); err != nil {
			m.logger.Errorw("Failed to recover worker",
				logx.Field("worker_id", w.WorkerID),
				logx.Field("error", err))
			continue
		}
		recovered++
	}

	m.logger.Infow("Workers recovered from database",
		logx.Field("total", len(workers)),
		logx.Field("recovered", recovered))

	return nil
}

// asyncUpdateWorker 异步批量更新Worker到数据库
func (m *WorkerManager) asyncUpdateWorker() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.UpdateInterval)
	defer ticker.Stop()

	batch := make([]*WorkerInfo, 0, m.config.UpdateBatchSize)

	for {
		select {
		case <-m.ctx.Done():
			// 刷新剩余批次
			if len(batch) > 0 {
				m.flushWorkerBatch(batch)
			}
			return

		case workerInfo := <-m.updateQueue:
			batch = append(batch, workerInfo)

			// 达到批次大小，立即刷新
			if len(batch) >= m.config.UpdateBatchSize {
				m.flushWorkerBatch(batch)
				batch = make([]*WorkerInfo, 0, m.config.UpdateBatchSize)
			}

		case <-ticker.C:
			// 定时刷新
			if len(batch) > 0 {
				m.flushWorkerBatch(batch)
				batch = make([]*WorkerInfo, 0, m.config.UpdateBatchSize)
			}
		}
	}
}

// flushWorkerBatch 刷新Worker批次到数据库
func (m *WorkerManager) flushWorkerBatch(batch []*WorkerInfo) {
	if len(batch) == 0 {
		return
	}

	ctx := context.Background()
	systemCtx := hooks.NewSystemContext(ctx)

	updated := 0
	created := 0

	for _, info := range batch {
		// 尝试更新
		err := m.db.Proxy.Update().
			Where(proxy.WorkerIDEQ(info.WorkerID)).
			SetWorkerStatus(proxy.WorkerStatus(info.Status)).
			SetLastHeartbeat(info.LastHeartbeat).
			SetCPUUsage(info.Load.CPUUsage).
			SetMemoryUsage(info.Load.MemoryUsage).
			SetDiskUsage(info.Load.DiskUsage).
			SetNetworkIn(info.Load.NetworkIn).
			SetNetworkOut(info.Load.NetworkOut).
			SetActiveSessions(info.Load.ActiveSessions).
			SetTotalRequests(info.Load.TotalRequests).
			SetSuccessCount(info.Load.SuccessCount).
			SetFailureCount(info.Load.FailureCount).
			Exec(systemCtx)

		if err != nil {
			if ent.IsNotFound(err) {
				// Worker不存在，创建新记录
				if err := m.createWorkerInDB(systemCtx, info); err != nil {
					m.logger.Errorw("Failed to create worker in DB",
						logx.Field("worker_id", info.WorkerID),
						logx.Field("error", err))
				} else {
					created++
				}
			} else {
				m.logger.Errorw("Failed to update worker in DB",
					logx.Field("worker_id", info.WorkerID),
					logx.Field("error", err))
			}
		} else {
			updated++
		}
	}

	m.logger.Infow("Flushed worker batch to database",
		logx.Field("batch_size", len(batch)),
		logx.Field("updated", updated),
		logx.Field("created", created))
}

// createWorkerInDB 在数据库中创建Worker记录
func (m *WorkerManager) createWorkerInDB(ctx context.Context, info *WorkerInfo) error {
	_, err := m.db.Proxy.Create().
		SetWorkerID(info.WorkerID).
		SetName(info.Name).
		SetIP(info.IP).
		SetPort(info.Port).
		SetNillableVersion(&info.Version).
		SetRegion(info.Region).
		SetZone(info.Zone).
		SetCapabilities(info.Capabilities).
		SetTags(info.Tags).
		SetEndpoints(info.Endpoints).
		SetWorkerStatus(proxy.WorkerStatus(info.Status)).
		SetLastHeartbeat(info.LastHeartbeat).
		SetRegisterTime(info.RegisterTime).
		SetWeight(info.Weight).
		SetPriority(info.Priority).
		SetCPUUsage(info.Load.CPUUsage).
		SetMemoryUsage(info.Load.MemoryUsage).
		SetDiskUsage(info.Load.DiskUsage).
		SetNetworkIn(info.Load.NetworkIn).
		SetNetworkOut(info.Load.NetworkOut).
		SetActiveSessions(info.Load.ActiveSessions).
		SetTotalRequests(info.Load.TotalRequests).
		SetSuccessCount(info.Load.SuccessCount).
		SetFailureCount(info.Load.FailureCount).
		SetMaxSessions(info.MaxSessions).
		SetHealthCheckFailures(info.HealthCheckFailures).
		SetHealthCheckURL(info.HealthCheckURL).
		SetLocalIP(info.LocalIP).
		SetPublicIP(info.PublicIP).
		SetNetworkSegments(info.NetworkSegments).
		SetMetadata(info.Metadata).
		SetTenantID(info.TenantID).
		SetStatus(1). // 启用状态
		Save(ctx)

	return err
}

// ttlCleanup TTL过期清理任务
func (m *WorkerManager) ttlCleanup() {
	defer m.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// 清理内存中的过期Worker
			expiredCount := m.registry.CleanExpired()

			if expiredCount > 0 {
				m.logger.Infow("TTL cleanup completed",
					logx.Field("expired_count", expiredCount))

				// 标记数据库中的Worker为offline
				m.markExpiredWorkersOffline()
			}
		}
	}
}

// markExpiredWorkersOffline 标记过期的Worker为offline
func (m *WorkerManager) markExpiredWorkersOffline() {
	ctx := context.Background()
	systemCtx := hooks.NewSystemContext(ctx)

	// 查询超时的Worker
	threshold := time.Now().Add(-m.config.TTL)

	err := m.db.Proxy.Update().
		Where(
			proxy.LastHeartbeatLT(threshold),
			proxy.WorkerStatusNEQ(proxy.WorkerStatusOffline),
		).
		SetWorkerStatus(proxy.WorkerStatusOffline).
		Exec(systemCtx)

	if err != nil {
		m.logger.Errorw("Failed to mark expired workers as offline",
			logx.Field("error", err))
	}
}

// asyncCollectMetrics 异步收集指标快照
func (m *WorkerManager) asyncCollectMetrics() {
	defer m.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectMetricsSnapshot()
		}
	}
}

// collectMetricsSnapshot 收集所有Worker的指标快照
func (m *WorkerManager) collectMetricsSnapshot() {
	ctx := context.Background()
	systemCtx := hooks.NewSystemContext(ctx)

	workers := m.registry.ListAll(false)
	now := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, w := range workers {
		// 获取上一次的指标快照
		prev, hasPrev := m.previousMetrics[w.WorkerID]

		// 计算增量值
		var networkInDelta, networkOutDelta uint64
		var requestDelta, successDelta, failureDelta uint64

		if hasPrev {
			// 计算增量（处理计数器溢出）
			if w.Load.NetworkIn >= prev.NetworkIn {
				networkInDelta = w.Load.NetworkIn - prev.NetworkIn
			}
			if w.Load.NetworkOut >= prev.NetworkOut {
				networkOutDelta = w.Load.NetworkOut - prev.NetworkOut
			}
			if w.Load.TotalRequests >= prev.TotalRequests {
				requestDelta = w.Load.TotalRequests - prev.TotalRequests
			}
			if w.Load.SuccessCount >= prev.SuccessCount {
				successDelta = w.Load.SuccessCount - prev.SuccessCount
			}
			if w.Load.FailureCount >= prev.FailureCount {
				failureDelta = w.Load.FailureCount - prev.FailureCount
			}
		}

		// 获取平均延迟（从最近的健康检查）
		avgLatency, _ := m.healthCheckLatency[w.WorkerID]

		// 保存指标快照
		_, err := m.db.ProxyMetrics.Create().
			SetWorkerID(w.WorkerID).
			SetTimestamp(now).
			SetCPUUsage(w.Load.CPUUsage).
			SetMemoryUsage(w.Load.MemoryUsage).
			SetDiskUsage(w.Load.DiskUsage).
			SetNetworkInDelta(networkInDelta).
			SetNetworkOutDelta(networkOutDelta).
			SetActiveSessions(w.Load.ActiveSessions).
			SetRequestCountDelta(requestDelta).
			SetSuccessCountDelta(successDelta).
			SetFailureCountDelta(failureDelta).
			SetAvgLatencyMs(avgLatency).
			SetWorkerStatus(w.Status).
			SetTenantID(w.TenantID).
			Save(systemCtx)

		if err != nil {
			m.logger.Errorw("Failed to save worker metrics",
				logx.Field("worker_id", w.WorkerID),
				logx.Field("error", err))
		}

		// 更新上一次指标快照
		m.previousMetrics[w.WorkerID] = &WorkerMetricsSnapshot{
			Timestamp:     now,
			NetworkIn:     w.Load.NetworkIn,
			NetworkOut:    w.Load.NetworkOut,
			TotalRequests: w.Load.TotalRequests,
			SuccessCount:  w.Load.SuccessCount,
			FailureCount:  w.Load.FailureCount,
		}
	}

	m.logger.Infow("Metrics snapshot collected",
		logx.Field("worker_count", len(workers)))
}

// healthCheckLoop 健康检查循环
func (m *WorkerManager) healthCheckLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.performHealthChecks()
		}
	}
}

// performHealthChecks 执行健康检查
func (m *WorkerManager) performHealthChecks() {
	workers := m.registry.ListAll(false)

	if len(workers) == 0 {
		return
	}

	m.logger.Infow("Starting health check cycle",
		logx.Field("worker_count", len(workers)))

	// 使用批量详细检查
	stats, results := m.healthChecker.BatchDetailedCheck(m.ctx, workers)

	// 更新Worker状态和延迟信息
	m.mu.Lock()
	statusChanged := 0
	for _, result := range results {
		// 记录健康检查延迟
		if result.Success {
			m.healthCheckLatency[result.WorkerID] = float64(result.LatencyMs)
		}

		// 更新内存中的健康检查状态
		err := m.registry.UpdateHealthCheck(result.WorkerID, result.Success)
		if err != nil {
			m.logger.Errorw("Failed to update health check status",
				logx.Field("worker_id", result.WorkerID),
				logx.Field("error", err))
			continue
		}

		// 记录状态变化
		if result.PreviousStatus != result.NewStatus {
			statusChanged++
			m.logger.Infow("Worker status changed",
				logx.Field("worker_id", result.WorkerID),
				logx.Field("previous_status", result.PreviousStatus),
				logx.Field("new_status", result.NewStatus),
				logx.Field("failure_count", getWorkerFailureCount(result.WorkerID, workers)))
		}
	}
	m.mu.Unlock()

	// 记录统计信息
	m.logger.Infow("Health check cycle completed",
		logx.Field("total_checked", stats.TotalChecked),
		logx.Field("success", stats.SuccessCount),
		logx.Field("failure", stats.FailureCount),
		logx.Field("avg_latency_ms", stats.AvgLatencyMs),
		logx.Field("status_changed", statusChanged))

	// 记录状态变化详情
	if len(stats.StatusChanges) > 0 {
		for change, count := range stats.StatusChanges {
			m.logger.Infow("Status change summary",
				logx.Field("change", change),
				logx.Field("count", count))
		}
	}
}

// getWorkerFailureCount 获取Worker的失败次数（辅助函数）
func getWorkerFailureCount(workerID string, workers []*WorkerInfo) int {
	for _, w := range workers {
		if w.WorkerID == workerID {
			return w.HealthCheckFailures
		}
	}
	return 0
}

// metricsCleanup 清理旧指标数据
func (m *WorkerManager) metricsCleanup() {
	defer m.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.cleanOldMetrics()
		}
	}
}

// cleanOldMetrics 清理旧指标
func (m *WorkerManager) cleanOldMetrics() {
	ctx := context.Background()
	systemCtx := hooks.NewSystemContext(ctx)

	threshold := time.Now().Add(-m.config.MetricsRetention)

	deleted, err := m.db.ProxyMetrics.Delete().
		Where(proxymetrics.TimestampLT(threshold)).
		Exec(systemCtx)

	if err != nil {
		m.logger.Errorw("Failed to clean old metrics",
			logx.Field("error", err))
	} else if deleted > 0 {
		m.logger.Infow("Old metrics cleaned",
			logx.Field("deleted_count", deleted),
			logx.Field("threshold", threshold))
	}
}

// convertEntToProxyInfo 将Ent Proxy实体转换为WorkerInfo (保持向后兼容)
func (m *WorkerManager) convertEntToProxyInfo(p *ent.Proxy) *WorkerInfo {
	// 处理指针类型的时间字段
	lastHeartbeat := time.Time{}
	if p.LastHeartbeat != nil {
		lastHeartbeat = *p.LastHeartbeat
	}

	registerTime := time.Time{}
	if p.RegisterTime != nil {
		registerTime = *p.RegisterTime
	}

	lastHealthCheck := time.Time{}
	if p.LastHealthCheck != nil {
		lastHealthCheck = *p.LastHealthCheck
	}

	return &WorkerInfo{
		WorkerID:            p.WorkerID,            // 映射：proxy_id → WorkerID
		Name:                p.Name,
		IP:                  p.IP,
		Port:                p.Port,
		Version:             p.Version,
		Region:              p.Region,
		Zone:                p.Zone,
		Capabilities:        p.Capabilities,
		Tags:                p.Tags,
		Endpoints:           p.Endpoints,
		Status:              string(p.WorkerStatus), // 映射：proxy_status → Status
		LastHeartbeat:       lastHeartbeat,
		RegisterTime:        registerTime,
		Weight:              p.Weight,
		Priority:            p.Priority,
		Load: WorkerLoad{
			CPUUsage:       p.CPUUsage,
			MemoryUsage:    p.MemoryUsage,
			DiskUsage:      p.DiskUsage,
			NetworkIn:      p.NetworkIn,
			NetworkOut:     p.NetworkOut,
			ActiveSessions: p.ActiveSessions,
			TotalRequests:  p.TotalRequests,
			SuccessCount:   p.SuccessCount,
			FailureCount:   p.FailureCount,
		},
		HealthCheckFailures: p.HealthCheckFailures,
		HealthCheckURL:      p.HealthCheckURL,
		LastHealthCheck:     lastHealthCheck,
		LocalIP:             p.LocalIP,
		PublicIP:            p.PublicIP,
		NetworkSegments:     p.NetworkSegments,
		TenantID:            p.TenantID,
		Metadata:            p.Metadata,
		MaxSessions:         p.MaxSessions,
	}
}

// convertEntToWorkerInfo 别名，保持向后兼容
func (m *WorkerManager) convertEntToWorkerInfo(p *ent.Proxy) *WorkerInfo {
	return m.convertEntToProxyInfo(p)
}

// ==================== 指标查询方法 ====================

// WorkerMetricsQueryOptions 指标查询选项
type WorkerMetricsQueryOptions struct {
	WorkerID  string     // 必须：Worker ID
	TenantID  uint64     // 必须：租户ID
	StartTime *time.Time // 可选：开始时间
	EndTime   *time.Time // 可选：结束时间
	Limit     int        // 可选：返回记录数量限制（默认100，最大1000）
}

// WorkerMetricsResult 指标查询结果
type WorkerMetricsResult struct {
	WorkerID         string    `json:"worker_id"`
	Timestamp        time.Time `json:"timestamp"`
	CPUUsage         float64   `json:"cpu_usage"`
	MemoryUsage      float64   `json:"memory_usage"`
	DiskUsage        float64   `json:"disk_usage"`
	NetworkInDelta   uint64    `json:"network_in_delta"`
	NetworkOutDelta  uint64    `json:"network_out_delta"`
	ActiveSessions   int       `json:"active_sessions"`
	RequestDelta     uint64    `json:"request_delta"`
	SuccessDelta     uint64    `json:"success_delta"`
	FailureDelta     uint64    `json:"failure_delta"`
	AvgLatencyMs     float64   `json:"avg_latency_ms"`
	WorkerStatus     string    `json:"worker_status"`
}

// GetWorkerMetrics 查询Worker的历史指标
func (m *WorkerManager) GetWorkerMetrics(ctx context.Context, opts *WorkerMetricsQueryOptions) ([]*WorkerMetricsResult, error) {
	if opts == nil || opts.WorkerID == "" {
		return nil, fmt.Errorf("worker_id is required")
	}

	// 设置默认limit
	limit := 100
	if opts.Limit > 0 {
		limit = opts.Limit
		if limit > 1000 {
			limit = 1000 // 最大1000条
		}
	}

	// 构建查询
	query := m.db.ProxyMetrics.Query().
		Where(
			proxymetrics.WorkerIDEQ(opts.WorkerID),
			proxymetrics.TenantIDEQ(opts.TenantID),
		)

	// 添加时间范围过滤
	if opts.StartTime != nil {
		query = query.Where(proxymetrics.TimestampGTE(*opts.StartTime))
	}
	if opts.EndTime != nil {
		query = query.Where(proxymetrics.TimestampLTE(*opts.EndTime))
	}

	// 执行查询（按时间降序）
	metrics, err := query.
		Order(ent.Desc(proxymetrics.FieldTimestamp)).
		Limit(limit).
		All(ctx)

	if err != nil {
		m.logger.Errorw("Failed to query worker metrics",
			logx.Field("worker_id", opts.WorkerID),
			logx.Field("tenant_id", opts.TenantID),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to query worker metrics: %w", err)
	}

	// 转换结果
	results := make([]*WorkerMetricsResult, 0, len(metrics))
	for _, metric := range metrics {
		results = append(results, &WorkerMetricsResult{
			WorkerID:        metric.WorkerID,
			Timestamp:       metric.Timestamp,
			CPUUsage:        metric.CPUUsage,
			MemoryUsage:     metric.MemoryUsage,
			DiskUsage:       metric.DiskUsage,
			NetworkInDelta:  metric.NetworkInDelta,
			NetworkOutDelta: metric.NetworkOutDelta,
			ActiveSessions:  metric.ActiveSessions,
			RequestDelta:    metric.RequestCountDelta,
			SuccessDelta:    metric.SuccessCountDelta,
			FailureDelta:    metric.FailureCountDelta,
			AvgLatencyMs:    metric.AvgLatencyMs,
			WorkerStatus:    metric.WorkerStatus,
		})
	}

	return results, nil
}

// WorkerMetricsStats 指标统计信息
type WorkerMetricsStats struct {
	WorkerID        string    `json:"worker_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	AvgCPU          float64   `json:"avg_cpu"`
	MaxCPU          float64   `json:"max_cpu"`
	AvgMemory       float64   `json:"avg_memory"`
	MaxMemory       float64   `json:"max_memory"`
	AvgSessions     float64   `json:"avg_sessions"`
	MaxSessions     int       `json:"max_sessions"`
	TotalRequests   uint64    `json:"total_requests"`
	TotalSuccess    uint64    `json:"total_success"`
	TotalFailure    uint64    `json:"total_failure"`
	SuccessRate     float64   `json:"success_rate"` // 成功率百分比
	AvgLatency      float64   `json:"avg_latency"`  // 平均延迟(ms)
	DataPoints      int       `json:"data_points"`  // 数据点数量
}

// GetWorkerMetricsStats 获取Worker的指标统计
func (m *WorkerManager) GetWorkerMetricsStats(ctx context.Context, opts *WorkerMetricsQueryOptions) (*WorkerMetricsStats, error) {
	metrics, err := m.GetWorkerMetrics(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics found for worker %s", opts.WorkerID)
	}

	// 计算统计信息
	stats := &WorkerMetricsStats{
		WorkerID:   opts.WorkerID,
		StartTime:  metrics[len(metrics)-1].Timestamp, // 最早的
		EndTime:    metrics[0].Timestamp,              // 最新的
		DataPoints: len(metrics),
	}

	var sumCPU, sumMemory, sumSessions, sumLatency float64
	var maxCPU, maxMemory float64
	var maxSessions int
	var totalRequests, totalSuccess, totalFailure uint64
	var latencyCount int

	for _, m := range metrics {
		// CPU统计
		sumCPU += m.CPUUsage
		if m.CPUUsage > maxCPU {
			maxCPU = m.CPUUsage
		}

		// 内存统计
		sumMemory += m.MemoryUsage
		if m.MemoryUsage > maxMemory {
			maxMemory = m.MemoryUsage
		}

		// 会话统计
		sumSessions += float64(m.ActiveSessions)
		if m.ActiveSessions > maxSessions {
			maxSessions = m.ActiveSessions
		}

		// 请求统计
		totalRequests += m.RequestDelta
		totalSuccess += m.SuccessDelta
		totalFailure += m.FailureDelta

		// 延迟统计
		if m.AvgLatencyMs > 0 {
			sumLatency += float64(m.AvgLatencyMs)
			latencyCount++
		}
	}

	// 计算平均值
	count := float64(len(metrics))
	stats.AvgCPU = sumCPU / count
	stats.MaxCPU = maxCPU
	stats.AvgMemory = sumMemory / count
	stats.MaxMemory = maxMemory
	stats.AvgSessions = sumSessions / count
	stats.MaxSessions = maxSessions
	stats.TotalRequests = totalRequests
	stats.TotalSuccess = totalSuccess
	stats.TotalFailure = totalFailure

	// 计算成功率
	if totalRequests > 0 {
		stats.SuccessRate = float64(totalSuccess) / float64(totalRequests) * 100
	}

	// 计算平均延迟
	if latencyCount > 0 {
		stats.AvgLatency = sumLatency / float64(latencyCount)
	}

	return stats, nil
}

// GetAllWorkersMetricsSummary 获取所有Worker的当前指标汇总
func (m *WorkerManager) GetAllWorkersMetricsSummary(tenantID uint64) map[string]interface{} {
	workers := m.registry.List(tenantID, false)

	totalWorkers := len(workers)
	onlineCount := 0
	degradedCount := 0
	offlineCount := 0
	totalSessions := 0
	totalRequests := uint64(0)
	totalSuccess := uint64(0)
	totalFailure := uint64(0)
	avgCPU := 0.0
	avgMemory := 0.0

	for _, w := range workers {
		switch w.Status {
		case "online":
			onlineCount++
		case "degraded":
			degradedCount++
		case "offline":
			offlineCount++
		}

		totalSessions += w.Load.ActiveSessions
		totalRequests += w.Load.TotalRequests
		totalSuccess += w.Load.SuccessCount
		totalFailure += w.Load.FailureCount
		avgCPU += w.Load.CPUUsage
		avgMemory += w.Load.MemoryUsage
	}

	if totalWorkers > 0 {
		avgCPU /= float64(totalWorkers)
		avgMemory /= float64(totalWorkers)
	}

	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(totalSuccess) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"total_workers":   totalWorkers,
		"online_count":    onlineCount,
		"degraded_count":  degradedCount,
		"offline_count":   offlineCount,
		"total_sessions":  totalSessions,
		"total_requests":  totalRequests,
		"total_success":   totalSuccess,
		"total_failure":   totalFailure,
		"success_rate":    successRate,
		"avg_cpu_usage":   avgCPU,
		"avg_memory_usage": avgMemory,
	}
}
