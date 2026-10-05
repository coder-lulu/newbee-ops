package proxygroup

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy分组列表（需要JWT）
func NewGetProxyGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupListLogic {
	return &GetProxyGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyGroupListLogic) GetProxyGroupList(req *types.ProxyGroupListReq) (resp *types.ProxyGroupListResp, err error) {
	rpcReq := &pb.ProxyGroupListReq{Page: req.Page, PageSize: req.PageSize}
	if req.Name != "" {
		rpcReq.Name = &req.Name
	}
	if req.Status != 0 {
		rpcReq.Status = &req.Status
	}
	result, err := l.svcCtx.OpsClient.GetProxyGroupList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	items := make([]types.ProxyGroupItem, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, types.ProxyGroupItem{Id: item.GetId(),
			Name:                item.GetName(),
			Description:         item.GetDescription(),
			SelectionStrategy:   item.GetSelectionStrategy(),
			HealthCheckInterval: int(item.GetHealthCheckInterval()),
			AutoFailover:        item.GetAutoFailover(),
			MemberCount:         int(item.GetTotalMembers()),
			OnlineCount:         int(item.GetOnlineMembers()),
			TotalWeight:         int(item.GetTotalWeight()),
			Status:              item.GetStatus()})
	}
	return &types.ProxyGroupListResp{Msg: "success", Data: types.ProxyGroupListInfo{Total: result.Total, Data: items}}, nil
}
