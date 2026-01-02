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

type CreateAgentGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentGroupMemberLogic {
	return &CreateAgentGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentGroupMemberLogic) CreateAgentGroupMember(in *ops.AgentGroupMemberInfo) (*ops.BaseIDResp, error) {
    query := l.svcCtx.DB.AgentGroupMember.Create().
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

	result, err := query.Save(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess }, nil
}
