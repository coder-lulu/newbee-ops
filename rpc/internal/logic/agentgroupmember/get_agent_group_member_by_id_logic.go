package agentgroupmember

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupMemberByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentGroupMemberByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupMemberByIdLogic {
	return &GetAgentGroupMemberByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentGroupMemberByIdLogic) GetAgentGroupMemberById(in *ops.IDReq) (*ops.AgentGroupMemberInfo, error) {
	result, err := l.svcCtx.DB.AgentGroupMember.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.AgentGroupMemberInfo{
		Id:          &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:	pointy.GetPointer(uint32(result.Status)),
		AgentId:	&result.AgentID,
		AgentGroupId:	&result.AgentGroupID,
		Priority:	pointy.GetPointer(int64(result.Priority)),
		Weight:	pointy.GetPointer(int64(result.Weight)),
		JoinedAt:	pointy.GetUnixMilliPointer(result.JoinedAt.UnixMilli()),
	}, nil
}

