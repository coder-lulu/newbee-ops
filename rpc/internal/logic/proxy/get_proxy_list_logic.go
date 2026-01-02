package proxy

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetProxyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyListLogic {
	return &GetProxyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyListLogic) GetProxyList(in *ops.ProxyListReq) (*ops.ProxyListResp, error) {
	var predicates []predicate.Proxy
	if in.CreatedAt != nil {
		predicates = append(predicates, proxy.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, proxy.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, proxy.StatusEQ(uint8(*in.Status)))
	}
	if in.ProxyId != nil {
		predicates = append(predicates, proxy.WorkerIDContains(*in.ProxyId))
	}
	if in.Name != nil {
		predicates = append(predicates, proxy.NameContains(*in.Name))
	}
	if in.Ip != nil {
		predicates = append(predicates, proxy.IPContains(*in.Ip))
	}
	if in.Port != nil {
		predicates = append(predicates, proxy.PortEQ(int(*in.Port)))
	}
	if in.Version != nil {
		predicates = append(predicates, proxy.VersionContains(*in.Version))
	}
	if in.Region != nil {
		predicates = append(predicates, proxy.RegionContains(*in.Region))
	}
	if in.Zone != nil {
		predicates = append(predicates, proxy.ZoneContains(*in.Zone))
	}
	if in.Capabilities != nil {
	}
	if in.Tags != nil {
	}
	if in.Endpoints != nil {
	}
	if in.ProxyStatus != nil {
		predicates = append(predicates, proxy.WorkerStatusEQ(stringToProxyStatus(in.ProxyStatus)))
	}
	if in.LastHeartbeat != nil {
		predicates = append(predicates, proxy.LastHeartbeatGTE(time.UnixMilli(*in.LastHeartbeat)))
	}
	if in.RegisterTime != nil {
		predicates = append(predicates, proxy.RegisterTimeGTE(time.UnixMilli(*in.RegisterTime)))
	}
	if in.Weight != nil {
		predicates = append(predicates, proxy.WeightEQ(int(*in.Weight)))
	}
	if in.Priority != nil {
		predicates = append(predicates, proxy.PriorityEQ(int(*in.Priority)))
	}
	if in.CpuUsage != nil {
		predicates = append(predicates, proxy.CPUUsageEQ(*in.CpuUsage))
	}
	if in.MemoryUsage != nil {
		predicates = append(predicates, proxy.MemoryUsageEQ(*in.MemoryUsage))
	}
	if in.DiskUsage != nil {
		predicates = append(predicates, proxy.DiskUsageEQ(*in.DiskUsage))
	}
	if in.NetworkIn != nil {
		predicates = append(predicates, proxy.NetworkInEQ(*in.NetworkIn))
	}
	if in.NetworkOut != nil {
		predicates = append(predicates, proxy.NetworkOutEQ(*in.NetworkOut))
	}
	if in.ActiveSessions != nil {
		predicates = append(predicates, proxy.ActiveSessionsEQ(int(*in.ActiveSessions)))
	}
	if in.TotalRequests != nil {
		predicates = append(predicates, proxy.TotalRequestsEQ(*in.TotalRequests))
	}
	if in.SuccessCount != nil {
		predicates = append(predicates, proxy.SuccessCountEQ(*in.SuccessCount))
	}
	if in.FailureCount != nil {
		predicates = append(predicates, proxy.FailureCountEQ(*in.FailureCount))
	}
	if in.MaxSessions != nil {
		predicates = append(predicates, proxy.MaxSessionsEQ(int(*in.MaxSessions)))
	}
	if in.HealthCheckFailures != nil {
		predicates = append(predicates, proxy.HealthCheckFailuresEQ(int(*in.HealthCheckFailures)))
	}
	if in.LastHealthCheck != nil {
		predicates = append(predicates, proxy.LastHealthCheckGTE(time.UnixMilli(*in.LastHealthCheck)))
	}
	if in.HealthCheckUrl != nil {
		predicates = append(predicates, proxy.HealthCheckURLContains(*in.HealthCheckUrl))
	}
	if in.LocalIp != nil {
		predicates = append(predicates, proxy.LocalIPContains(*in.LocalIp))
	}
	if in.PublicIp != nil {
		predicates = append(predicates, proxy.PublicIPContains(*in.PublicIp))
	}
	if in.NetworkSegments != nil {
	}
	if in.Metadata != nil {
	}
	if in.LastError != nil {
		predicates = append(predicates, proxy.LastErrorContains(*in.LastError))
	}
	result, err := l.svcCtx.DB.Proxy.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ProxyListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.ProxyInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:	pointy.GetPointer(uint32(v.Status)),
			ProxyId:	&v.WorkerID,
			Name:	&v.Name,
			Ip:	&v.IP,
			Port:	pointy.GetPointer(int64(v.Port)),
			Version:	&v.Version,
			Region:	&v.Region,
			Zone:	&v.Zone,
			Capabilities: pointy.GetPointer(mustMarshalJSON(v.Capabilities)),
			Tags: pointy.GetPointer(mustMarshalJSON(v.Tags)),
			Endpoints: pointy.GetPointer(mustMarshalJSON(v.Endpoints)),
			ProxyStatus: pointy.GetPointer(proxyStatusToString(v.WorkerStatus)),
			LastHeartbeat:	pointy.GetUnixMilliPointer(v.LastHeartbeat.UnixMilli()),
			RegisterTime:	pointy.GetUnixMilliPointer(v.RegisterTime.UnixMilli()),
			Weight:	pointy.GetPointer(int64(v.Weight)),
			Priority:	pointy.GetPointer(int64(v.Priority)),
			CpuUsage:	&v.CPUUsage,
			MemoryUsage:	&v.MemoryUsage,
			DiskUsage:	&v.DiskUsage,
			NetworkIn:	&v.NetworkIn,
			NetworkOut:	&v.NetworkOut,
			ActiveSessions:	pointy.GetPointer(int64(v.ActiveSessions)),
			TotalRequests:	&v.TotalRequests,
			SuccessCount:	&v.SuccessCount,
			FailureCount:	&v.FailureCount,
			MaxSessions:	pointy.GetPointer(int64(v.MaxSessions)),
			HealthCheckFailures:	pointy.GetPointer(int64(v.HealthCheckFailures)),
			LastHealthCheck:	pointy.GetUnixMilliPointer(v.LastHealthCheck.UnixMilli()),
			HealthCheckUrl:	&v.HealthCheckURL,
			LocalIp:	&v.LocalIP,
			PublicIp:	&v.PublicIP,
			NetworkSegments: pointy.GetPointer(mustMarshalJSON(v.NetworkSegments)),
			Metadata: pointy.GetPointer(mustMarshalJSON(v.Metadata)),
			LastError:	&v.LastError,
		})
	}

	return resp, nil
}
