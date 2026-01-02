package proxy_group

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProxyGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyGroupLogic {
	return &UpdateProxyGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProxyGroupLogic) UpdateProxyGroup(in *ops.ProxyGroupInfo) (*ops.BaseResp, error) {
	query:= l.svcCtx.DB.ProxyGroup.UpdateOneID(*in.Id).
			SetNotNilDepartmentID(in.DepartmentId).
			SetNotNilName(in.Name).
			SetNotNilDescription(in.Description).
			SetSelectionStrategy(stringToSelectionStrategy(in.SelectionStrategy)).
			SetNotNilAutoFailover(in.AutoFailover)

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.HealthCheckInterval != nil {
		query.SetNotNilHealthCheckInterval(pointy.GetPointer(int(*in.HealthCheckInterval)))
	}
	if in.MaxRetryCount != nil {
		query.SetNotNilMaxRetryCount(pointy.GetPointer(int(*in.MaxRetryCount)))
	}
	if in.MinHealthyProxies != nil {
		query.SetNotNilMinHealthyWorkers(pointy.GetPointer(int(*in.MinHealthyProxies)))
	}
	if in.ConnectionTimeout != nil {
		query.SetNotNilConnectionTimeout(pointy.GetPointer(int(*in.ConnectionTimeout)))
	}
	if in.RequestTimeout != nil {
		query.SetNotNilRequestTimeout(pointy.GetPointer(int(*in.RequestTimeout)))
	}
	if in.TotalMembers != nil {
		query.SetNotNilTotalMembers(pointy.GetPointer(int(*in.TotalMembers)))
	}
	if in.OnlineMembers != nil {
		query.SetNotNilOnlineMembers(pointy.GetPointer(int(*in.OnlineMembers)))
	}
	if in.TotalWeight != nil {
		query.SetNotNilTotalWeight(pointy.GetPointer(int(*in.TotalWeight)))
	}

	 err := query.Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
