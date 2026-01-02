package worker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// HealthChecker Worker健康检查器
type HealthChecker struct {
	client  *http.Client
	timeout time.Duration
	logger  logx.Logger
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(timeout time.Duration, logger logx.Logger) *HealthChecker {
	return &HealthChecker{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		timeout: timeout,
		logger:  logger,
	}
}

// CheckWorker 执行单个Worker的健康检查
func (h *HealthChecker) CheckWorker(ctx context.Context, worker *WorkerInfo) bool {
	if worker.HealthCheckURL == "" {
		// 没有配置健康检查URL，认为健康（依赖心跳）
		return true
	}

	// 创建带超时的请求
	req, err := http.NewRequestWithContext(ctx, "GET", worker.HealthCheckURL, nil)
	if err != nil {
		h.logger.Errorw("Failed to create health check request",
			logx.Field("worker_id", worker.WorkerID),
			logx.Field("url", worker.HealthCheckURL),
			logx.Field("error", err))
		return false
	}

	// 添加基本的请求头
	req.Header.Set("User-Agent", "Ops-Center-HealthChecker/1.0")
	req.Header.Set("Accept", "*/*")

	// 执行请求
	startTime := time.Now()
	resp, err := h.client.Do(req)
	latency := time.Since(startTime)

	if err != nil {
		h.logger.Errorw("Health check request failed",
			logx.Field("worker_id", worker.WorkerID),
			logx.Field("url", worker.HealthCheckURL),
			logx.Field("latency_ms", latency.Milliseconds()),
			logx.Field("error", err))
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	if success {
		h.logger.Infow("Health check passed",
			logx.Field("worker_id", worker.WorkerID),
			logx.Field("status_code", resp.StatusCode),
			logx.Field("latency_ms", latency.Milliseconds()))
	} else {
		h.logger.Errorw("Health check failed - bad status code",
			logx.Field("worker_id", worker.WorkerID),
			logx.Field("status_code", resp.StatusCode),
			logx.Field("latency_ms", latency.Milliseconds()))
	}

	return success
}

// CheckWorkers 批量检查Worker健康状态（并发执行）
func (h *HealthChecker) CheckWorkers(ctx context.Context, workers []*WorkerInfo) map[string]bool {
	results := make(map[string]bool)
	resultsMu := sync.Mutex{}

	// 使用WaitGroup等待所有检查完成
	var wg sync.WaitGroup

	// 使用信号量限制并发数（最多20个并发检查）
	sem := make(chan struct{}, 20)

	for _, worker := range workers {
		// 跳过没有配置健康检查URL的Worker
		if worker.HealthCheckURL == "" {
			resultsMu.Lock()
			results[worker.WorkerID] = true // 默认健康
			resultsMu.Unlock()
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(w *WorkerInfo) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			// 执行健康检查
			success := h.CheckWorker(ctx, w)

			// 记录结果
			resultsMu.Lock()
			results[w.WorkerID] = success
			resultsMu.Unlock()
		}(worker)
	}

	// 等待所有检查完成
	wg.Wait()

	return results
}

// CheckWorkerWithRetry 带重试的健康检查
func (h *HealthChecker) CheckWorkerWithRetry(ctx context.Context, worker *WorkerInfo, maxRetries int) bool {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前短暂等待
			select {
			case <-ctx.Done():
				return false
			case <-time.After(time.Second):
				// 继续
			}

			h.logger.Infow("Retrying health check",
				logx.Field("worker_id", worker.WorkerID),
				logx.Field("attempt", attempt+1),
				logx.Field("max_retries", maxRetries))
		}

		if h.CheckWorker(ctx, worker) {
			return true
		}
	}

	return false
}

// HealthCheckResult 健康检查结果
type HealthCheckResult struct {
	WorkerID    string
	Success     bool
	StatusCode  int
	LatencyMs   int64
	Error       string
	CheckedAt   time.Time
	PreviousStatus string
	NewStatus   string
}

// DetailedCheck 执行详细的健康检查（返回详细结果）
func (h *HealthChecker) DetailedCheck(ctx context.Context, worker *WorkerInfo) *HealthCheckResult {
	result := &HealthCheckResult{
		WorkerID:       worker.WorkerID,
		CheckedAt:      time.Now(),
		PreviousStatus: worker.Status,
	}

	if worker.HealthCheckURL == "" {
		result.Success = true
		result.NewStatus = worker.Status
		return result
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", worker.HealthCheckURL, nil)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		result.NewStatus = calculateNewStatus(worker.Status, worker.HealthCheckFailures+1)
		return result
	}

	req.Header.Set("User-Agent", "Ops-Center-HealthChecker/1.0")

	// 执行请求
	startTime := time.Now()
	resp, err := h.client.Do(req)
	result.LatencyMs = time.Since(startTime).Milliseconds()

	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("request failed: %v", err)
		result.NewStatus = calculateNewStatus(worker.Status, worker.HealthCheckFailures+1)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	if result.Success {
		// 健康检查成功 - 状态可能恢复
		result.NewStatus = "online"
	} else {
		// 健康检查失败 - 状态可能降级
		result.Error = fmt.Sprintf("bad status code: %d", resp.StatusCode)
		result.NewStatus = calculateNewStatus(worker.Status, worker.HealthCheckFailures+1)
	}

	return result
}

// calculateNewStatus 根据当前状态和失败次数计算新状态
func calculateNewStatus(currentStatus string, failureCount int) string {
	if failureCount >= 5 {
		return "offline"
	} else if failureCount >= 3 {
		return "degraded"
	}
	return currentStatus
}

// HealthCheckStats 健康检查统计
type HealthCheckStats struct {
	TotalChecked   int
	SuccessCount   int
	FailureCount   int
	SkippedCount   int
	AvgLatencyMs   int64
	StatusChanges  map[string]int // online->degraded, degraded->offline等
}

// BatchDetailedCheck 批量详细健康检查（返回统计信息）
func (h *HealthChecker) BatchDetailedCheck(ctx context.Context, workers []*WorkerInfo) (*HealthCheckStats, []*HealthCheckResult) {
	stats := &HealthCheckStats{
		StatusChanges: make(map[string]int),
	}

	results := make([]*HealthCheckResult, 0, len(workers))
	resultsChan := make(chan *HealthCheckResult, len(workers))

	// 使用WaitGroup和信号量
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for _, worker := range workers {
		wg.Add(1)
		sem <- struct{}{}

		go func(w *WorkerInfo) {
			defer wg.Done()
			defer func() { <-sem }()

			result := h.DetailedCheck(ctx, w)
			resultsChan <- result
		}(worker)
	}

	// 等待所有检查完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果和统计
	var totalLatency int64
	for result := range resultsChan {
		results = append(results, result)
		stats.TotalChecked++

		if result.Success {
			stats.SuccessCount++
		} else if result.Error != "" {
			stats.FailureCount++
		} else {
			stats.SkippedCount++
		}

		totalLatency += result.LatencyMs

		// 记录状态变化
		if result.PreviousStatus != result.NewStatus {
			key := fmt.Sprintf("%s->%s", result.PreviousStatus, result.NewStatus)
			stats.StatusChanges[key]++
		}
	}

	if stats.TotalChecked > 0 {
		stats.AvgLatencyMs = totalLatency / int64(stats.TotalChecked)
	}

	return stats, results
}

// IsHealthy 快速检查Worker是否健康（不执行HTTP请求，仅基于状态）
func IsHealthy(worker *WorkerInfo) bool {
	return worker.Status == "online" || worker.Status == "degraded"
}

// ShouldPerformCheck 判断是否应该执行健康检查
func ShouldPerformCheck(worker *WorkerInfo, lastCheckTime time.Time, checkInterval time.Duration) bool {
	// 如果没有配置健康检查URL，不需要检查
	if worker.HealthCheckURL == "" {
		return false
	}

	// 如果从未检查过，应该检查
	if lastCheckTime.IsZero() {
		return true
	}

	// 如果距离上次检查超过间隔时间，应该检查
	return time.Since(lastCheckTime) >= checkInterval
}
