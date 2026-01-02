package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyLogic {
	return &UpdateProxyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProxyLogic) UpdateProxy(in *ops.ProxyInfo) (*ops.BaseResp, error) {
	// Convert JSON strings to native types
	capabilities := unmarshalJSONStringSlice(in.Capabilities)
	tags := unmarshalJSONStringSlice(in.Tags)
	endpoints := unmarshalJSONStringMap(in.Endpoints)
	proxyStatus := stringToProxyStatus(in.ProxyStatus)
	networkSegments := unmarshalJSONStringSlice(in.NetworkSegments)
	metadata := unmarshalJSONInterfaceMap(in.Metadata)
	query:= l.svcCtx.DB.Proxy.UpdateOneID(*in.Id).
			SetNotNilWorkerID(in.ProxyId).
			SetNotNilName(in.Name).
			SetNotNilIP(in.Ip).
			SetNotNilVersion(in.Version).
			SetNotNilRegion(in.Region).
			SetNotNilZone(in.Zone).
			SetCapabilities(capabilities).
			SetTags(tags).
			SetEndpoints(endpoints).
			SetWorkerStatus(proxyStatus).
			SetNotNilLastHeartbeat(pointy.GetTimeMilliPointer(in.LastHeartbeat)).
			SetNotNilRegisterTime(pointy.GetTimeMilliPointer(in.RegisterTime)).
			SetNotNilCPUUsage(in.CpuUsage).
			SetNotNilMemoryUsage(in.MemoryUsage).
			SetNotNilDiskUsage(in.DiskUsage).
			SetNotNilNetworkIn(in.NetworkIn).
			SetNotNilNetworkOut(in.NetworkOut).
			SetNotNilTotalRequests(in.TotalRequests).
			SetNotNilSuccessCount(in.SuccessCount).
			SetNotNilFailureCount(in.FailureCount).
			SetNotNilLastHealthCheck(pointy.GetTimeMilliPointer(in.LastHealthCheck)).
			SetNotNilHealthCheckURL(in.HealthCheckUrl).
			SetNotNilLocalIP(in.LocalIp).
			SetNotNilPublicIP(in.PublicIp).
			SetNetworkSegments(networkSegments).
			SetMetadata(metadata).
			SetNotNilLastError(in.LastError)

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.Port != nil {
		query.SetNotNilPort(pointy.GetPointer(int(*in.Port)))
	}
	if in.Weight != nil {
		query.SetNotNilWeight(pointy.GetPointer(int(*in.Weight)))
	}
	if in.Priority != nil {
		query.SetNotNilPriority(pointy.GetPointer(int(*in.Priority)))
	}
	if in.ActiveSessions != nil {
		query.SetNotNilActiveSessions(pointy.GetPointer(int(*in.ActiveSessions)))
	}
	if in.MaxSessions != nil {
		query.SetNotNilMaxSessions(pointy.GetPointer(int(*in.MaxSessions)))
	}
	if in.HealthCheckFailures != nil {
		query.SetNotNilHealthCheckFailures(pointy.GetPointer(int(*in.HealthCheckFailures)))
	}

	 err := query.Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
