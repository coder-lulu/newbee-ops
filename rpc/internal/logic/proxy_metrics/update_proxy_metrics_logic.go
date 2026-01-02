package proxy_metrics

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyMetricsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProxyMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyMetricsLogic {
	return &UpdateProxyMetricsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProxyMetricsLogic) UpdateProxyMetrics(in *ops.ProxyMetricsInfo) (*ops.BaseResp, error) {
	query:= l.svcCtx.DB.ProxyMetrics.UpdateOneID(*in.Id).
			SetNotNilWorkerID(in.ProxyId).
			SetNotNilTimestamp(pointy.GetTimeMilliPointer(in.Timestamp)).
			SetNotNilCPUUsage(in.CpuUsage).
			SetNotNilMemoryUsage(in.MemoryUsage).
			SetNotNilDiskUsage(in.DiskUsage).
			SetNotNilNetworkInDelta(in.NetworkInDelta).
			SetNotNilNetworkOutDelta(in.NetworkOutDelta).
			SetNotNilRequestCountDelta(in.RequestCountDelta).
			SetNotNilSuccessCountDelta(in.SuccessCountDelta).
			SetNotNilFailureCountDelta(in.FailureCountDelta).
			SetNotNilAvgLatencyMs(in.AvgLatencyMs).
			SetNotNilWorkerStatus(in.ProxyStatus)

	if in.ActiveSessions != nil {
		query.SetNotNilActiveSessions(pointy.GetPointer(int(*in.ActiveSessions)))
	}

	 err := query.Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
