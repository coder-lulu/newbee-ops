package proxy_metrics

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyMetricsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyMetricsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyMetricsByIdLogic {
	return &GetProxyMetricsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyMetricsByIdLogic) GetProxyMetricsById(in *ops.IDReq) (*ops.ProxyMetricsInfo, error) {
	result, err := l.svcCtx.DB.ProxyMetrics.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ProxyMetricsInfo{
		Id:          &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		ProxyId:	&result.WorkerID,
		Timestamp:	pointy.GetPointer(result.Timestamp.UnixMilli()),
		CpuUsage:	&result.CPUUsage,
		MemoryUsage:	&result.MemoryUsage,
		DiskUsage:	&result.DiskUsage,
		NetworkInDelta:	&result.NetworkInDelta,
		NetworkOutDelta:	&result.NetworkOutDelta,
		ActiveSessions:	pointy.GetPointer(int64(result.ActiveSessions)),
		RequestCountDelta:	&result.RequestCountDelta,
		SuccessCountDelta:	&result.SuccessCountDelta,
		FailureCountDelta:	&result.FailureCountDelta,
		AvgLatencyMs:	&result.AvgLatencyMs,
		ProxyStatus:	&result.WorkerStatus,
	}, nil
}

