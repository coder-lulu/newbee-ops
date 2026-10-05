package proxy

import (
	"context"
	"encoding/json"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy详情（需要JWT）
func NewGetProxyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyByIdLogic {
	return &GetProxyByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyByIdLogic) GetProxyById(req *types.IDReq) (resp *types.ProxyDetail, err error) {
	item, err := l.svcCtx.OpsClient.GetProxyById(l.ctx, &pb.IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	var capabilities, tags, segments []string
	_ = json.Unmarshal([]byte(item.GetCapabilities()), &capabilities)
	_ = json.Unmarshal([]byte(item.GetTags()), &tags)
	_ = json.Unmarshal([]byte(item.GetNetworkSegments()), &segments)
	return &types.ProxyDetail{Capabilities: capabilities, Tags: tags, NetworkSegments: segments, Id: item.GetId(),
		CreatedAt:           item.GetCreatedAt(),
		UpdatedAt:           item.GetUpdatedAt(),
		Status:              item.GetStatus(),
		WorkerID:            item.GetProxyId(),
		Name:                item.GetName(),
		IP:                  item.GetIp(),
		Port:                int(item.GetPort()),
		Version:             item.GetVersion(),
		Region:              item.GetRegion(),
		Zone:                item.GetZone(),
		Endpoints:           item.GetEndpoints(),
		WorkerStatus:        item.GetProxyStatus(),
		LastHeartbeat:       item.GetLastHeartbeat(),
		RegisterTime:        item.GetRegisterTime(),
		Weight:              int(item.GetWeight()),
		Priority:            int(item.GetPriority()),
		CPUUsage:            item.GetCpuUsage(),
		MemoryUsage:         item.GetMemoryUsage(),
		DiskUsage:           item.GetDiskUsage(),
		NetworkIn:           item.GetNetworkIn(),
		NetworkOut:          item.GetNetworkOut(),
		ActiveSessions:      int(item.GetActiveSessions()),
		TotalRequests:       item.GetTotalRequests(),
		SuccessCount:        item.GetSuccessCount(),
		FailureCount:        item.GetFailureCount(),
		MaxSessions:         int(item.GetMaxSessions()),
		HealthCheckFailures: int(item.GetHealthCheckFailures()),
		HealthCheckURL:      item.GetHealthCheckUrl(),
		LastHealthCheck:     item.GetLastHealthCheck(),
		LocalIP:             item.GetLocalIp(),
		PublicIP:            item.GetPublicIp(),
		Metadata:            item.GetMetadata()}, nil
}
