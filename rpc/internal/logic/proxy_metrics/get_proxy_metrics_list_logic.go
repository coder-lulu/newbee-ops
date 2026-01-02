package proxy_metrics

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/proxymetrics"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetProxyMetricsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyMetricsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyMetricsListLogic {
	return &GetProxyMetricsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyMetricsListLogic) GetProxyMetricsList(in *ops.ProxyMetricsListReq) (*ops.ProxyMetricsListResp, error) {
	var predicates []predicate.ProxyMetrics
	if in.CreatedAt != nil {
		predicates = append(predicates, proxymetrics.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, proxymetrics.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.ProxyId != nil {
		predicates = append(predicates, proxymetrics.WorkerIDContains(*in.ProxyId))
	}
	if in.Timestamp != nil {
		predicates = append(predicates, proxymetrics.TimestampGTE(time.UnixMilli(*in.Timestamp)))
	}
	if in.CpuUsage != nil {
		predicates = append(predicates, proxymetrics.CPUUsageEQ(*in.CpuUsage))
	}
	if in.MemoryUsage != nil {
		predicates = append(predicates, proxymetrics.MemoryUsageEQ(*in.MemoryUsage))
	}
	if in.DiskUsage != nil {
		predicates = append(predicates, proxymetrics.DiskUsageEQ(*in.DiskUsage))
	}
	if in.NetworkInDelta != nil {
		predicates = append(predicates, proxymetrics.NetworkInDeltaEQ(*in.NetworkInDelta))
	}
	if in.NetworkOutDelta != nil {
		predicates = append(predicates, proxymetrics.NetworkOutDeltaEQ(*in.NetworkOutDelta))
	}
	if in.ActiveSessions != nil {
		predicates = append(predicates, proxymetrics.ActiveSessionsEQ(int(*in.ActiveSessions)))
	}
	if in.RequestCountDelta != nil {
		predicates = append(predicates, proxymetrics.RequestCountDeltaEQ(*in.RequestCountDelta))
	}
	if in.SuccessCountDelta != nil {
		predicates = append(predicates, proxymetrics.SuccessCountDeltaEQ(*in.SuccessCountDelta))
	}
	if in.FailureCountDelta != nil {
		predicates = append(predicates, proxymetrics.FailureCountDeltaEQ(*in.FailureCountDelta))
	}
	if in.AvgLatencyMs != nil {
		predicates = append(predicates, proxymetrics.AvgLatencyMsEQ(*in.AvgLatencyMs))
	}
	if in.ProxyStatus != nil {
		predicates = append(predicates, proxymetrics.WorkerStatusContains(*in.ProxyStatus))
	}
	result, err := l.svcCtx.DB.ProxyMetrics.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ProxyMetricsListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.ProxyMetricsInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			ProxyId:	&v.WorkerID,
			Timestamp:	pointy.GetPointer(v.Timestamp.UnixMilli()),
			CpuUsage:	&v.CPUUsage,
			MemoryUsage:	&v.MemoryUsage,
			DiskUsage:	&v.DiskUsage,
			NetworkInDelta:	&v.NetworkInDelta,
			NetworkOutDelta:	&v.NetworkOutDelta,
			ActiveSessions:	pointy.GetPointer(int64(v.ActiveSessions)),
			RequestCountDelta:	&v.RequestCountDelta,
			SuccessCountDelta:	&v.SuccessCountDelta,
			FailureCountDelta:	&v.FailureCountDelta,
			AvgLatencyMs:	&v.AvgLatencyMs,
			ProxyStatus:	&v.WorkerStatus,
		})
	}

	return resp, nil
}
