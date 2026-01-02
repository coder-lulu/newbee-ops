package worker

import (
	"context"
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strings"
	"sync/atomic"
)

// ========== 1. Round Robin Strategy (轮询) ==========

type RoundRobinStrategy struct {
	counter uint64
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{}
}

func (s *RoundRobinStrategy) Name() string {
	return "round_robin"
}

func (s *RoundRobinStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 轮询选择
	idx := atomic.AddUint64(&s.counter, 1) % uint64(len(filtered))
	return filtered[idx], nil
}

// ========== 2. Least Connections Strategy (最少连接 - 默认) ==========

type LeastConnectionsStrategy struct{}

func NewLeastConnectionsStrategy() *LeastConnectionsStrategy {
	return &LeastConnectionsStrategy{}
}

func (s *LeastConnectionsStrategy) Name() string {
	return "least_connections"
}

func (s *LeastConnectionsStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 按活跃会话数升序排序
	sort.Slice(filtered, func(i, j int) bool {
		// 优先选择online状态
		if filtered[i].Status != filtered[j].Status {
			return filtered[i].Status == "online"
		}
		// 相同状态，选择会话数少的
		return filtered[i].Load.ActiveSessions < filtered[j].Load.ActiveSessions
	})

	return filtered[0], nil
}

// ========== 3. Weighted Round Robin Strategy (加权轮询) ==========

type WeightedRoundRobinStrategy struct {
	counter uint64
}

func NewWeightedRoundRobinStrategy() *WeightedRoundRobinStrategy {
	return &WeightedRoundRobinStrategy{}
}

func (s *WeightedRoundRobinStrategy) Name() string {
	return "weighted_round_robin"
}

func (s *WeightedRoundRobinStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 构建加权列表
	weighted := make([]*WorkerInfo, 0)
	for _, w := range filtered {
		weight := w.Weight

		// 根据状态调整权重
		switch w.Status {
		case "online":
			// 保持原权重
		case "degraded":
			// 降级状态权重减半
			weight = weight / 2
		case "offline":
			// 离线状态权重为0（实际上已被过滤）
			weight = 0
		}

		// 根据权重添加到列表（权重越大，出现次数越多）
		for i := 0; i < weight; i++ {
			weighted = append(weighted, w)
		}
	}

	if len(weighted) == 0 {
		// 所有Worker权重都为0，降级为普通轮询
		idx := atomic.AddUint64(&s.counter, 1) % uint64(len(filtered))
		return filtered[idx], nil
	}

	// 轮询选择
	idx := atomic.AddUint64(&s.counter, 1) % uint64(len(weighted))
	return weighted[idx], nil
}

// ========== 4. Consistent Hash Strategy (一致性哈希 - 会话亲和性) ==========

type ConsistentHashStrategy struct{}

func NewConsistentHashStrategy() *ConsistentHashStrategy {
	return &ConsistentHashStrategy{}
}

func (s *ConsistentHashStrategy) Name() string {
	return "consistent_hash"
}

func (s *ConsistentHashStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 检查是否有SessionID
	if opts == nil || opts.SessionID == "" {
		// 无SessionID，降级为最少连接策略
		return NewLeastConnectionsStrategy().Select(ctx, filtered, opts)
	}

	// 使用CRC32哈希
	hash := crc32.ChecksumIEEE([]byte(opts.SessionID))
	idx := int(hash) % len(filtered)

	return filtered[idx], nil
}

// ========== 5. Geo Nearest Strategy (地理位置优先) ==========

type GeoNearestStrategy struct{}

func NewGeoNearestStrategy() *GeoNearestStrategy {
	return &GeoNearestStrategy{}
}

func (s *GeoNearestStrategy) Name() string {
	return "geo_nearest"
}

func (s *GeoNearestStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 如果指定了优先区域
	if opts != nil && opts.PreferredRegion != "" {
		// 先找完全匹配Region和Zone的
		if opts.PreferredZone != "" {
			for _, w := range filtered {
				if w.Region == opts.PreferredRegion && w.Zone == opts.PreferredZone {
					return w, nil
				}
			}
		}

		// 再找匹配Region的
		regionMatched := make([]*WorkerInfo, 0)
		for _, w := range filtered {
			if w.Region == opts.PreferredRegion {
				regionMatched = append(regionMatched, w)
			}
		}

		if len(regionMatched) > 0 {
			// 在匹配的Region中选择最少连接
			return NewLeastConnectionsStrategy().Select(ctx, regionMatched, nil)
		}
	}

	// 没有匹配的地理位置，降级为最少连接
	return NewLeastConnectionsStrategy().Select(ctx, filtered, nil)
}

// ========== 6. Priority Strategy (优先级) ==========

type PriorityStrategy struct{}

func NewPriorityStrategy() *PriorityStrategy {
	return &PriorityStrategy{}
}

func (s *PriorityStrategy) Name() string {
	return "priority"
}

func (s *PriorityStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 按优先级降序排序
	sort.Slice(filtered, func(i, j int) bool {
		// 优先按优先级排序
		if filtered[i].Priority != filtered[j].Priority {
			return filtered[i].Priority > filtered[j].Priority
		}
		// 相同优先级，按状态排序
		if filtered[i].Status != filtered[j].Status {
			return filtered[i].Status == "online"
		}
		// 相同优先级和状态，按活跃会话数排序
		return filtered[i].Load.ActiveSessions < filtered[j].Load.ActiveSessions
	})

	return filtered[0], nil
}

// ========== 7. Random Strategy (随机 - 额外的简单策略) ==========

type RandomStrategy struct{}

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{}
}

func (s *RandomStrategy) Name() string {
	return "random"
}

func (s *RandomStrategy) Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error) {
	if len(workers) == 0 {
		return nil, fmt.Errorf("no available workers")
	}

	// 过滤Worker
	filtered := filterWorkers(workers, opts)
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no workers match the selection criteria")
	}

	// 随机选择
	idx := rand.Intn(len(filtered))
	return filtered[idx], nil
}

// ========== 过滤逻辑 ==========

// filterWorkers 根据选项过滤Worker列表
func filterWorkers(workers []*WorkerInfo, opts *SelectOptions) []*WorkerInfo {
	if opts == nil {
		// 无过滤条件，只排除offline
		return filterByStatus(workers)
	}

	filtered := workers

	// 1. 过滤状态 (排除offline)
	filtered = filterByStatus(filtered)

	// 2. 过滤必须的能力
	if len(opts.RequiredCapabilities) > 0 {
		filtered = filterByCapabilities(filtered, opts.RequiredCapabilities)
	}

	// 3. 过滤必须的标签
	if len(opts.RequiredTags) > 0 {
		filtered = filterByTags(filtered, opts.RequiredTags)
	}

	// 4. 过滤最小权重
	if opts.MinWeight > 0 {
		filtered = filterByMinWeight(filtered, opts.MinWeight)
	}

	// 5. 排除指定的Worker ID
	if len(opts.ExcludeWorkerIDs) > 0 {
		filtered = filterByExclude(filtered, opts.ExcludeWorkerIDs)
	}

	return filtered
}

// filterByStatus 过滤状态（排除offline）
func filterByStatus(workers []*WorkerInfo) []*WorkerInfo {
	result := make([]*WorkerInfo, 0, len(workers))
	for _, w := range workers {
		if w.Status != "offline" {
			result = append(result, w)
		}
	}
	return result
}

// filterByCapabilities 过滤能力
func filterByCapabilities(workers []*WorkerInfo, required []string) []*WorkerInfo {
	result := make([]*WorkerInfo, 0, len(workers))
	for _, w := range workers {
		if hasAllCapabilities(w.Capabilities, required) {
			result = append(result, w)
		}
	}
	return result
}

// hasAllCapabilities 检查是否包含所有必须的能力
func hasAllCapabilities(capabilities []string, required []string) bool {
	capSet := make(map[string]bool)
	for _, cap := range capabilities {
		capSet[strings.ToLower(cap)] = true
	}

	for _, req := range required {
		if !capSet[strings.ToLower(req)] {
			return false
		}
	}
	return true
}

// filterByTags 过滤标签
func filterByTags(workers []*WorkerInfo, required []string) []*WorkerInfo {
	result := make([]*WorkerInfo, 0, len(workers))
	for _, w := range workers {
		if hasAllTags(w.Tags, required) {
			result = append(result, w)
		}
	}
	return result
}

// hasAllTags 检查是否包含所有必须的标签
func hasAllTags(tags []string, required []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[strings.ToLower(tag)] = true
	}

	for _, req := range required {
		if !tagSet[strings.ToLower(req)] {
			return false
		}
	}
	return true
}

// filterByMinWeight 过滤最小权重
func filterByMinWeight(workers []*WorkerInfo, minWeight int) []*WorkerInfo {
	result := make([]*WorkerInfo, 0, len(workers))
	for _, w := range workers {
		if w.Weight >= minWeight {
			result = append(result, w)
		}
	}
	return result
}

// filterByExclude 排除指定的Worker ID
func filterByExclude(workers []*WorkerInfo, exclude []string) []*WorkerInfo {
	excludeSet := make(map[string]bool)
	for _, id := range exclude {
		excludeSet[id] = true
	}

	result := make([]*WorkerInfo, 0, len(workers))
	for _, w := range workers {
		if !excludeSet[w.WorkerID] {
			result = append(result, w)
		}
	}
	return result
}
