package agentgroupmember

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveAgentsFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveAgentsFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveAgentsFromGroupLogic {
	return &RemoveAgentsFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RemoveAgentsFromGroup 从组中移除Agent
func (l *RemoveAgentsFromGroupLogic) RemoveAgentsFromGroup(in *ops.AgentGroupMemberReq) (*ops.BaseResp, error) {
	if in.AgentGroupId == nil {
		return nil, fmt.Errorf("agent_group_id is required")
	}

	if len(in.AgentIds) == 0 {
		return nil, fmt.Errorf("at least one agent_id is required")
	}

	// 批量删除成员
	_, err := l.svcCtx.DB.AgentGroupMember.Delete().
		Where(
			agentgroupmember.AgentGroupIDEQ(*in.AgentGroupId),
			agentgroupmember.AgentIDIn(in.AgentIds...),
		).
		Exec(l.ctx)

	if err != nil {
		return nil, err
	}

	return &ops.BaseResp{
		Msg: "Agents removed from group successfully",
	}, nil
}
