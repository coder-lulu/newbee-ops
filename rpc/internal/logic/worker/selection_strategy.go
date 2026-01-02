package worker

import "context"

// SelectionStrategy Worker选择策略接口
type SelectionStrategy interface {
	// Name 返回策略名称
	Name() string

	// Select 从候选Worker列表中选择一个Worker
	Select(ctx context.Context, workers []*WorkerInfo, opts *SelectOptions) (*WorkerInfo, error)
}

// SelectOptions Worker选择选项
type SelectOptions struct {
	// 必须的能力
	RequiredCapabilities []string

	// 优先区域
	PreferredRegion string

	// 优先可用区
	PreferredZone string

	// 必须的标签
	RequiredTags []string

	// 最小权重
	MinWeight int

	// 会话亲和性（用于一致性哈希）
	SessionID string

	// 排除的Worker ID列表
	ExcludeWorkerIDs []string
}
