package proxygroup

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy分组详情（需要JWT）
func NewGetProxyGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupByIdLogic {
	return &GetProxyGroupByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyGroupByIdLogic) GetProxyGroupById(req *types.IDPathReq) (resp *types.ProxyGroupDetailResp, err error) {
	item, err := l.svcCtx.OpsClient.GetProxyGroupById(l.ctx, &pb.IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.ProxyGroupDetailResp{Msg: "success", Data: types.ProxyGroupDetail{Id: item.GetId(),
		CreatedAt:           item.GetCreatedAt(),
		UpdatedAt:           item.GetUpdatedAt(),
		Status:              item.GetStatus(),
		Name:                item.GetName(),
		Description:         item.GetDescription(),
		SelectionStrategy:   item.GetSelectionStrategy(),
		HealthCheckInterval: int(item.GetHealthCheckInterval()),
		AutoFailover:        item.GetAutoFailover(),
		MaxRetryCount:       int(item.GetMaxRetryCount()),
		MinHealthyWorkers:   int(item.GetMinHealthyProxies())}}, nil
}
