package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyByIdLogic {
	return &GetProxyByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyByIdLogic) GetProxyById(in *ops.IDReq) (*ops.ProxyInfo, error) {
	result, err := l.svcCtx.DB.Proxy.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// Convert native Go types to JSON strings for protobuf
	capabilitiesJSON := mustMarshalJSON(result.Capabilities)
	tagsJSON := mustMarshalJSON(result.Tags)
	endpointsJSON := mustMarshalJSON(result.Endpoints)
	networkSegmentsJSON := mustMarshalJSON(result.NetworkSegments)
	metadataJSON := mustMarshalJSON(result.Metadata)
	proxyStatusStr := proxyStatusToString(result.WorkerStatus)

	return &ops.ProxyInfo{
		Id:                  &result.ID,
		CreatedAt:           pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:           pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:              pointy.GetPointer(uint32(result.Status)),
		ProxyId:             &result.WorkerID,
		Name:                &result.Name,
		Ip:                  &result.IP,
		Port:                pointy.GetPointer(int64(result.Port)),
		Version:             &result.Version,
		Region:              &result.Region,
		Zone:                &result.Zone,
		Capabilities:        &capabilitiesJSON,
		Tags:                &tagsJSON,
		Endpoints:           &endpointsJSON,
		ProxyStatus:         &proxyStatusStr,
		LastHeartbeat:       optionalUnixMilli(result.LastHeartbeat),
		RegisterTime:        optionalUnixMilli(result.RegisterTime),
		Weight:              pointy.GetPointer(int64(result.Weight)),
		Priority:            pointy.GetPointer(int64(result.Priority)),
		CpuUsage:            &result.CPUUsage,
		MemoryUsage:         &result.MemoryUsage,
		DiskUsage:           &result.DiskUsage,
		NetworkIn:           &result.NetworkIn,
		NetworkOut:          &result.NetworkOut,
		ActiveSessions:      pointy.GetPointer(int64(result.ActiveSessions)),
		TotalRequests:       &result.TotalRequests,
		SuccessCount:        &result.SuccessCount,
		FailureCount:        &result.FailureCount,
		MaxSessions:         pointy.GetPointer(int64(result.MaxSessions)),
		HealthCheckFailures: pointy.GetPointer(int64(result.HealthCheckFailures)),
		LastHealthCheck:     optionalUnixMilli(result.LastHealthCheck),
		HealthCheckUrl:      &result.HealthCheckURL,
		LocalIp:             &result.LocalIP,
		PublicIp:            &result.PublicIP,
		NetworkSegments:     &networkSegmentsJSON,
		Metadata:            &metadataJSON,
		LastError:           &result.LastError,
	}, nil
}
