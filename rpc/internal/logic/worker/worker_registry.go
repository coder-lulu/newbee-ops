package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// WorkerLoad 负载信息
type WorkerLoad struct {
	CPUUsage       float64 // CPU使用率（0-100）
	MemoryUsage    float64 // 内存使用率（0-100）
	DiskUsage      float64 // 磁盘使用率（0-100）
	NetworkIn      uint64  // 网络入流量（字节）
	NetworkOut     uint64  // 网络出流量（字节）
	ActiveSessions int     // 当前活跃会话数
	TotalRequests  uint64  // 总请求次数
	SuccessCount   uint64  // 成功请求次数
	FailureCount   uint64  // 失败请求次数
}

// WorkerInfo Worker在内存中的完整信息
type WorkerInfo struct {
	// 基础信息
	WorkerID string
	Name     string
	IP       string
	Port     int
	Version  string

	// 地理位置
	Region string
	Zone   string

	// 能力与标签
	Capabilities []string
	Tags         []string
	Endpoints    map[string]string

	// 状态信息
	Status        string    // online/degraded/offline
	LastHeartbeat time.Time
	RegisterTime  time.Time

	// 负载均衡配置
	Weight   int
	Priority int

	// 负载信息
	Load WorkerLoad

	// 健康检查
	HealthCheckFailures int
	HealthCheckURL      string
	LastHealthCheck     time.Time

	// 网络信息
	LocalIP         string
	PublicIP        string
	NetworkSegments []string

	// 租户隔离
	TenantID uint64

	// 扩展元数据
	Metadata map[string]interface{}

	// 容量限制
	MaxSessions int
}

// WorkerRegistry Worker内存注册表
type WorkerRegistry struct {
	mu      sync.RWMutex
	workers map[string]*WorkerInfo // key: worker_id
	ttl     time.Duration

	logger logx.Logger
}

// NewWorkerRegistry 创建Worker注册表
func NewWorkerRegistry(ttl time.Duration, logger logx.Logger) *WorkerRegistry {
	return &WorkerRegistry{
		workers: make(map[string]*WorkerInfo),
		ttl:     ttl,
		logger:  logger,
	}
}

// Register 注册或更新Worker
func (r *WorkerRegistry) Register(info *WorkerInfo) error {
	if info == nil {
		return fmt.Errorf("worker info cannot be nil")
	}

	if info.WorkerID == "" {
		return fmt.Errorf("worker_id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// 检查是否已存在
	existing, exists := r.workers[info.WorkerID]
	if exists {
		// 更新已有Worker - 保留一些累计字段
		info.RegisterTime = existing.RegisterTime
		info.Load.TotalRequests = existing.Load.TotalRequests
		info.Load.SuccessCount = existing.Load.SuccessCount
		info.Load.FailureCount = existing.Load.FailureCount

		r.logger.Infow("Worker re-registered",
			logx.Field("worker_id", info.WorkerID),
			logx.Field("tenant_id", info.TenantID),
			logx.Field("ip", info.IP))
	} else {
		// 新注册
		info.RegisterTime = now
		r.logger.Infow("Worker registered",
			logx.Field("worker_id", info.WorkerID),
			logx.Field("tenant_id", info.TenantID),
			logx.Field("ip", info.IP))
	}

	// 更新心跳时间
	info.LastHeartbeat = now

	// 存储到内存
	r.workers[info.WorkerID] = info

	return nil
}

// Heartbeat 更新Worker心跳和状态
func (r *WorkerRegistry) Heartbeat(workerID string, tenantID uint64, status string, load *WorkerLoad) error {
	if workerID == "" {
		return fmt.Errorf("worker_id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	worker, exists := r.workers[workerID]
	if !exists {
		return fmt.Errorf("worker not found: %s", workerID)
	}

	// 验证租户
	if worker.TenantID != tenantID {
		return fmt.Errorf("tenant mismatch: worker belongs to tenant %d, got %d", worker.TenantID, tenantID)
	}

	// 更新心跳时间
	worker.LastHeartbeat = time.Now()

	// 更新状态
	if status != "" {
		worker.Status = status
	}

	// 更新负载信息
	if load != nil {
		worker.Load = *load
	}

	return nil
}

// Get 获取单个Worker
func (r *WorkerRegistry) Get(workerID string) (*WorkerInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	worker, exists := r.workers[workerID]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if r.isExpired(worker) {
		return nil, false
	}

	return worker, true
}

// List 获取Worker列表（按租户过滤）
func (r *WorkerRegistry) List(tenantID uint64, includeExpired bool) []*WorkerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*WorkerInfo

	for _, worker := range r.workers {
		// 租户过滤
		if worker.TenantID != tenantID {
			continue
		}

		// TTL过滤
		if !includeExpired && r.isExpired(worker) {
			continue
		}

		result = append(result, worker)
	}

	return result
}

// ListAll 获取所有Worker（不过滤租户，用于管理后台）
func (r *WorkerRegistry) ListAll(includeExpired bool) []*WorkerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*WorkerInfo

	for _, worker := range r.workers {
		// TTL过滤
		if !includeExpired && r.isExpired(worker) {
			continue
		}

		result = append(result, worker)
	}

	return result
}

// Remove 移除Worker
func (r *WorkerRegistry) Remove(workerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workers[workerID]; !exists {
		return fmt.Errorf("worker not found: %s", workerID)
	}

	delete(r.workers, workerID)

	r.logger.Infow("Worker removed",
		logx.Field("worker_id", workerID))

	return nil
}

// UpdateHealthCheck 更新健康检查状态
func (r *WorkerRegistry) UpdateHealthCheck(workerID string, success bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	worker, exists := r.workers[workerID]
	if !exists {
		return fmt.Errorf("worker not found: %s", workerID)
	}

	worker.LastHealthCheck = time.Now()

	if success {
		// 健康检查成功 - 清零失败计数，状态恢复为online
		worker.HealthCheckFailures = 0
		if worker.Status == "offline" || worker.Status == "degraded" {
			worker.Status = "online"
			r.logger.Infow("Worker health recovered",
				logx.Field("worker_id", workerID),
				logx.Field("tenant_id", worker.TenantID))
		}
	} else {
		// 健康检查失败 - 增加失败计数
		worker.HealthCheckFailures++

		// 状态转换逻辑
		if worker.HealthCheckFailures >= 5 {
			// 连续5次失败 → offline
			if worker.Status != "offline" {
				worker.Status = "offline"
				r.logger.Errorw("Worker marked as offline",
					logx.Field("worker_id", workerID),
					logx.Field("tenant_id", worker.TenantID),
					logx.Field("failures", worker.HealthCheckFailures))
			}
		} else if worker.HealthCheckFailures >= 3 {
			// 连续3次失败 → degraded
			if worker.Status == "online" {
				worker.Status = "degraded"
				r.logger.Errorw("Worker degraded",
					logx.Field("worker_id", workerID),
					logx.Field("tenant_id", worker.TenantID),
					logx.Field("failures", worker.HealthCheckFailures))
			}
		}
	}

	return nil
}

// CleanExpired 清理过期的Worker
func (r *WorkerRegistry) CleanExpired() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	expiredWorkers := []string{}

	for workerID, worker := range r.workers {
		if now.Sub(worker.LastHeartbeat) > r.ttl {
			expiredWorkers = append(expiredWorkers, workerID)
		}
	}

	// 删除过期的Worker
	for _, workerID := range expiredWorkers {
		worker := r.workers[workerID]
		delete(r.workers, workerID)

		r.logger.Infow("Worker expired and removed",
			logx.Field("worker_id", workerID),
			logx.Field("tenant_id", worker.TenantID),
			logx.Field("last_heartbeat", worker.LastHeartbeat))
	}

	return len(expiredWorkers)
}

// Count 获取Worker数量统计
func (r *WorkerRegistry) Count(tenantID uint64) map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := map[string]int{
		"total":    0,
		"online":   0,
		"degraded": 0,
		"offline":  0,
		"expired":  0,
	}

	for _, worker := range r.workers {
		// 租户过滤
		if tenantID > 0 && worker.TenantID != tenantID {
			continue
		}

		stats["total"]++

		if r.isExpired(worker) {
			stats["expired"]++
			continue
		}

		switch worker.Status {
		case "online":
			stats["online"]++
		case "degraded":
			stats["degraded"]++
		case "offline":
			stats["offline"]++
		}
	}

	return stats
}

// isExpired 检查Worker是否过期（内部方法，不加锁）
func (r *WorkerRegistry) isExpired(worker *WorkerInfo) bool {
	return time.Since(worker.LastHeartbeat) > r.ttl
}

// GetExpiredWorkers 获取过期的Worker列表（用于同步到数据库标记为offline）
func (r *WorkerRegistry) GetExpiredWorkers() []*WorkerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var expired []*WorkerInfo

	for _, worker := range r.workers {
		if r.isExpired(worker) {
			expired = append(expired, worker)
		}
	}

	return expired
}

// UpdateLoad 更新Worker负载信息（用于心跳）
func (r *WorkerRegistry) UpdateLoad(workerID string, load *WorkerLoad) error {
	if load == nil {
		return fmt.Errorf("load cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	worker, exists := r.workers[workerID]
	if !exists {
		return fmt.Errorf("worker not found: %s", workerID)
	}

	worker.Load = *load
	worker.LastHeartbeat = time.Now()

	return nil
}
