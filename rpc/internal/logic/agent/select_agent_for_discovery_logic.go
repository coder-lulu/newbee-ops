package agent

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/agent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectAgentForDiscoveryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectAgentForDiscoveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAgentForDiscoveryLogic {
	return &SelectAgentForDiscoveryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SelectAgentForDiscovery 为Discovery任务选择Agent
func (l *SelectAgentForDiscoveryLogic) SelectAgentForDiscovery(in *ops.AgentSelectionReq) (*ops.AgentSelectionResp, error) {
	// 根据选择模式执行不同的选择策略
	selectionMode := "auto_load_balance"
	if in.SelectionMode != nil {
		selectionMode = *in.SelectionMode
	}

	var selectedAgent *ent.Agent
	var err error

	switch selectionMode {
	case "manual_agent":
		selectedAgent, err = l.selectManualAgent(in)
	case "manual_group":
		selectedAgent, err = l.selectFromGroup(in)
	case "auto_load_balance":
		selectedAgent, err = l.selectByLoadBalance(in)
	case "auto_availability":
		selectedAgent, err = l.selectByAvailability(in)
	case "auto_performance":
		selectedAgent, err = l.selectByPerformance(in)
	case "auto_geo":
		selectedAgent, err = l.selectByGeography(in)
	default:
		return nil, fmt.Errorf("unsupported selection mode: %s", selectionMode)
	}

	if err != nil {
		return nil, err
	}

	if selectedAgent == nil {
		return nil, fmt.Errorf("no suitable agent found")
	}

	// 返回选中的Agent信息
	return &ops.AgentSelectionResp{
		AgentId: &selectedAgent.AgentID,
		Host:    &selectedAgent.Host,
		Port:    func() *int64 { p := int64(selectedAgent.Port); return &p }(),
		Msg:     func() *string { s := "Agent selected successfully"; return &s }(),
	}, nil
}

// selectManualAgent 手动指定Agent
func (l *SelectAgentForDiscoveryLogic) selectManualAgent(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	if in.AgentId == nil || *in.AgentId == "" {
		return nil, fmt.Errorf("agent_id is required for manual_agent mode")
	}

	// 查找指定的Agent
	selectedAgent, err := l.svcCtx.DB.Agent.Query().
		Where(agent.AgentIDEQ(*in.AgentId)).
		Only(l.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to find agent: %w", err)
	}

	return selectedAgent, nil
}

// selectFromGroup 从Agent组中选择
func (l *SelectAgentForDiscoveryLogic) selectFromGroup(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	if in.AgentGroupId == nil || *in.AgentGroupId == 0 {
		return nil, fmt.Errorf("agent_group_id is required for manual_group mode")
	}

	// 查找Agent组的所有成员
	members, err := l.svcCtx.DB.AgentGroupMember.Query().
		Where(agentgroupmember.AgentGroupIDEQ(*in.AgentGroupId)).
		WithAgent().
		All(l.ctx)

	if err != nil || len(members) == 0 {
		return nil, fmt.Errorf("no agents found in group: %w", err)
	}

	// 从组成员中选择第一个在线的Agent
	for _, member := range members {
		if member.Edges.Agent != nil && member.Edges.Agent.AgentStatus == "online" {
			return member.Edges.Agent, nil
		}
	}

	return nil, fmt.Errorf("no online agents in group")
}

// selectByLoadBalance 根据负载均衡选择
func (l *SelectAgentForDiscoveryLogic) selectByLoadBalance(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	candidates, err := l.getEligibleAgents(in)
	if err != nil || len(candidates) == 0 {
		return nil, fmt.Errorf("no eligible agents found: %w", err)
	}

	// 选择活跃会话数最少的Agent
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].ActiveSessions < candidates[j].ActiveSessions
	})

	return candidates[0], nil
}

// selectByAvailability 根据可用性选择
func (l *SelectAgentForDiscoveryLogic) selectByAvailability(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	candidates, err := l.getEligibleAgents(in)
	if err != nil || len(candidates) == 0 {
		return nil, fmt.Errorf("no eligible agents found: %w", err)
	}

	// 选择最近心跳的Agent
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].LastHeartbeat.After(*candidates[j].LastHeartbeat)
	})

	return candidates[0], nil
}

// selectByPerformance 根据性能选择
func (l *SelectAgentForDiscoveryLogic) selectByPerformance(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	candidates, err := l.getEligibleAgents(in)
	if err != nil || len(candidates) == 0 {
		return nil, fmt.Errorf("no eligible agents found: %w", err)
	}

	// 选择CPU和内存使用率最低的Agent
	sort.Slice(candidates, func(i, j int) bool {
		loadI := candidates[i].CPUUsage + candidates[i].MemoryUsage
		loadJ := candidates[j].CPUUsage + candidates[j].MemoryUsage
		return loadI < loadJ
	})

	return candidates[0], nil
}

// selectByGeography 根据地理位置选择
func (l *SelectAgentForDiscoveryLogic) selectByGeography(in *ops.AgentSelectionReq) (*ent.Agent, error) {
	candidates, err := l.getEligibleAgents(in)
	if err != nil || len(candidates) == 0 {
		return nil, fmt.Errorf("no eligible agents found: %w", err)
	}

	// 如果指定了region,优先选择同region的Agent
	if in.Region != nil && *in.Region != "" {
		for _, candidate := range candidates {
			if candidate.Region == *in.Region {
				return candidate, nil
			}
		}
	}

	// 否则随机选择一个
	if len(candidates) > 0 {
		rand.Seed(time.Now().UnixNano())
		return candidates[rand.Intn(len(candidates))], nil
	}

	return nil, fmt.Errorf("no agents in specified region")
}

// getEligibleAgents 获取符合条件的Agent列表
func (l *SelectAgentForDiscoveryLogic) getEligibleAgents(in *ops.AgentSelectionReq) ([]*ent.Agent, error) {
	query := l.svcCtx.DB.Agent.Query().
		Where(agent.AgentStatusEQ("online"))

	// 根据请求条件过滤
	if in.Region != nil && *in.Region != "" {
		query = query.Where(agent.RegionEQ(*in.Region))
	}

	agents, err := query.All(l.ctx)
	if err != nil {
		return nil, err
	}

	return agents, nil
}
