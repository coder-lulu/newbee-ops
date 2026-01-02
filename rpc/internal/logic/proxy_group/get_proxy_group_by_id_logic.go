package proxy_group

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupByIdLogic {
	return &GetProxyGroupByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyGroupByIdLogic) GetProxyGroupById(in *ops.IDReq) (*ops.ProxyGroupInfo, error) {
	result, err := l.svcCtx.DB.ProxyGroup.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ProxyGroupInfo{
		Id:          &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		DepartmentId:	&result.DepartmentID,
		Status:	pointy.GetPointer(uint32(result.Status)),
		Name:	&result.Name,
		Description:	&result.Description,
		SelectionStrategy: pointy.GetPointer(selectionStrategyToString(result.SelectionStrategy)),
		HealthCheckInterval:	pointy.GetPointer(int64(result.HealthCheckInterval)),
		AutoFailover:	&result.AutoFailover,
		MaxRetryCount:	pointy.GetPointer(int64(result.MaxRetryCount)),
		MinHealthyProxies:	pointy.GetPointer(int64(result.MinHealthyWorkers)),
		ConnectionTimeout:	pointy.GetPointer(int64(result.ConnectionTimeout)),
		RequestTimeout:	pointy.GetPointer(int64(result.RequestTimeout)),
		TotalMembers:	pointy.GetPointer(int64(result.TotalMembers)),
		OnlineMembers:	pointy.GetPointer(int64(result.OnlineMembers)),
		TotalWeight:	pointy.GetPointer(int64(result.TotalWeight)),
	}, nil
}

