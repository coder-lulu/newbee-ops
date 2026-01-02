package agent

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agent"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentHeartbeatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentHeartbeatLogic {
	return &UpdateAgentHeartbeatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateAgentHeartbeat 更新Agent心跳
func (l *UpdateAgentHeartbeatLogic) UpdateAgentHeartbeat(in *ops.AgentHeartbeatReq) (*ops.BaseResp, error) {
	// 根据agent_id查找Agent
	agentEntity, err := l.svcCtx.DB.Agent.Query().
		Where(agent.AgentIDEQ(*in.AgentId)).
		Only(l.ctx)

	if err != nil {
		return nil, err
	}

	// 更新心跳时间和状态
	now := time.Now()
	updateBuilder := l.svcCtx.DB.Agent.UpdateOneID(agentEntity.ID).
		SetLastHeartbeat(now)

	// 更新Agent状态
	if in.AgentStatus != nil {
		updateBuilder = updateBuilder.SetAgentStatus(agent.AgentStatus(*in.AgentStatus))
	}

	// 更新监控指标
	if in.ActiveSessions != nil {
		updateBuilder = updateBuilder.SetActiveSessions(int(*in.ActiveSessions))
	}
	if in.CpuUsage != nil {
		updateBuilder = updateBuilder.SetCPUUsage(*in.CpuUsage)
	}
	if in.MemoryUsage != nil {
		updateBuilder = updateBuilder.SetMemoryUsage(*in.MemoryUsage)
	}

	// 如果Agent状态是online，更新最后在线时间
	if in.AgentStatus != nil && *in.AgentStatus == "online" {
		updateBuilder = updateBuilder.SetLastOnlineAt(now)
	}

	// 执行更新
	err = updateBuilder.Exec(l.ctx)
	if err != nil {
		return nil, err
	}

	return &ops.BaseResp{
		Msg: "Agent heartbeat updated successfully",
	}, nil
}
