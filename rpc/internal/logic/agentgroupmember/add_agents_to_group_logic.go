package agentgroupmember

import (
	"context"
	"fmt"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAgentsToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAgentsToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAgentsToGroupLogic {
	return &AddAgentsToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddAgentsToGroup 添加Agent到组
func (l *AddAgentsToGroupLogic) AddAgentsToGroup(in *ops.AgentGroupMemberReq) (*ops.BaseResp, error) {
	if in.AgentGroupId == nil {
		return nil, fmt.Errorf("agent_group_id is required")
	}

	if len(in.AgentIds) == 0 {
		return nil, fmt.Errorf("at least one agent_id is required")
	}

	// 验证Agent Group存在
	_, err := l.svcCtx.DB.AgentGroup.Get(l.ctx, *in.AgentGroupId)
	if err != nil {
		return nil, err
	}

	// 批量添加成员
	now := time.Now()
	priority := 100
	if in.Priority != nil {
		priority = int(*in.Priority)
	}

	weight := 1
	if in.Weight != nil {
		weight = int(*in.Weight)
	}

	for _, agentID := range in.AgentIds {
		// 检查是否已经存在
		exists, err := l.svcCtx.DB.AgentGroupMember.Query().
			Where(
				agentgroupmember.AgentIDEQ(agentID),
				agentgroupmember.AgentGroupIDEQ(*in.AgentGroupId),
			).
			Exist(l.ctx)

		if err != nil {
			return nil, err
		}

		if exists {
			// 已存在，跳过
			continue
		}

		// 创建新成员
		_, err = l.svcCtx.DB.AgentGroupMember.Create().
			SetAgentID(agentID).
			SetAgentGroupID(*in.AgentGroupId).
			SetPriority(priority).
			SetWeight(weight).
			SetJoinedAt(now).
			Save(l.ctx)

		if err != nil {
			return nil, err
		}
	}

	return &ops.BaseResp{
		Msg: "Agents added to group successfully",
	}, nil
}
