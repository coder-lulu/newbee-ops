package proxy

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy列表（需要JWT）
func NewGetProxyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyListLogic {
	return &GetProxyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyListLogic) GetProxyList(req *types.ProxyListReq) (resp *types.ProxyListResp, err error) {
	rpcReq := &pb.ProxyListReq{Page: req.Page, PageSize: req.PageSize}
	if req.Name != "" {
		rpcReq.Name = &req.Name
	}
	if req.Status != 0 {
		rpcReq.Status = &req.Status
	}
	if req.ProxyStatus != "" {
		rpcReq.ProxyStatus = &req.ProxyStatus
	}
	if req.Region != "" {
		rpcReq.Region = &req.Region
	}
	if req.Zone != "" {
		rpcReq.Zone = &req.Zone
	}
	if req.Tags != "" {
		rpcReq.Tags = &req.Tags
	}
	result, err := l.svcCtx.OpsClient.GetProxyList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	items := make([]types.ProxyItem, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, types.ProxyItem{Id: item.GetId(),
			ProxyID:             item.GetProxyId(),
			Name:                item.GetName(),
			IP:                  item.GetIp(),
			Port:                int(item.GetPort()),
			Region:              item.GetRegion(),
			Zone:                item.GetZone(),
			ProxyStatus:         item.GetProxyStatus(),
			LastHeartbeat:       item.GetLastHeartbeat(),
			CPUUsage:            item.GetCpuUsage(),
			MemoryUsage:         item.GetMemoryUsage(),
			ActiveSessions:      int(item.GetActiveSessions()),
			MaxSessions:         int(item.GetMaxSessions()),
			TotalRequests:       item.GetTotalRequests(),
			SuccessCount:        item.GetSuccessCount(),
			Weight:              int(item.GetWeight()),
			Priority:            int(item.GetPriority()),
			HealthCheckFailures: int(item.GetHealthCheckFailures())})
	}
	return &types.ProxyListResp{Msg: "success", Data: types.ProxyListInfo{Total: result.Total, Data: items}}, nil
}
