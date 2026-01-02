package proxy_group_member

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupMemberByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyGroupMemberByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupMemberByIdLogic {
	return &GetProxyGroupMemberByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyGroupMemberByIdLogic) GetProxyGroupMemberById(in *ops.IDReq) (*ops.ProxyGroupMemberInfo, error) {
	result, err := l.svcCtx.DB.ProxyGroupMember.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ProxyGroupMemberInfo{
		Id:          &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		DepartmentId:	&result.DepartmentID,
		Status:	pointy.GetPointer(uint32(result.Status)),
		ProxyId:	&result.WorkerID,
		ProxyGroupId:	&result.WorkerGroupID,
		Weight:	pointy.GetPointer(int64(result.Weight)),
		Priority:	pointy.GetPointer(int64(result.Priority)),
		JoinedAt:	pointy.GetUnixMilliPointer(result.JoinedAt.UnixMilli()),
		LastSelectedAt:	pointy.GetUnixMilliPointer(result.LastSelectedAt.UnixMilli()),
		SelectCount:	pointy.GetPointer(int64(result.SelectCount)),
		IsBackup:	&result.IsBackup,
		ConsecutiveFailures:	pointy.GetPointer(int64(result.ConsecutiveFailures)),
		LastFailureAt:	pointy.GetUnixMilliPointer(result.LastFailureAt.UnixMilli()),
	}, nil
}

