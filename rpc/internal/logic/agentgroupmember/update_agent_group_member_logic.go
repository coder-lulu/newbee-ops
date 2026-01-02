package agentgroupmember

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentGroupMemberLogic {
	return &UpdateAgentGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAgentGroupMemberLogic) UpdateAgentGroupMember(in *ops.AgentGroupMemberInfo) (*ops.BaseResp, error) {
	query:= l.svcCtx.DB.AgentGroupMember.UpdateOneID(*in.Id).
			SetNotNilAgentID(in.AgentId).
			SetNotNilAgentGroupID(in.AgentGroupId).
			SetNotNilJoinedAt(pointy.GetTimeMilliPointer(in.JoinedAt))

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.Priority != nil {
		query.SetNotNilPriority(pointy.GetPointer(int(*in.Priority)))
	}
	if in.Weight != nil {
		query.SetNotNilWeight(pointy.GetPointer(int(*in.Weight)))
	}

	 err := query.Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
