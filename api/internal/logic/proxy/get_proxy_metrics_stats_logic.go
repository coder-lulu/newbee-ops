package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyMetricsStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询Proxy指标统计（需要JWT）
func NewGetProxyMetricsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyMetricsStatsLogic {
	return &GetProxyMetricsStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyMetricsStatsLogic) GetProxyMetricsStats(req *types.ProxyMetricsStatsReq) (resp *types.ProxyMetricsStatsResp, err error) {
	items, err := readMetrics(l.ctx, l.svcCtx.OpsClient, req.ProxyID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	stats := types.ProxyMetricsStatsData{ProxyID: req.ProxyID, StartTime: req.StartTime, EndTime: req.EndTime, DataPoints: len(items)}
	for _, x := range items {
		stats.AvgCPU += x.CPUUsage
		stats.AvgMemory += x.MemoryUsage
		stats.AvgSessions += float64(x.ActiveSessions)
		stats.AvgLatency += x.AvgLatencyMs
		stats.TotalRequests += x.RequestDelta
		stats.TotalSuccess += x.SuccessDelta
		stats.TotalFailure += x.FailureDelta
		if x.CPUUsage > stats.MaxCPU {
			stats.MaxCPU = x.CPUUsage
		}
		if x.MemoryUsage > stats.MaxMemory {
			stats.MaxMemory = x.MemoryUsage
		}
		if x.ActiveSessions > stats.MaxSessions {
			stats.MaxSessions = x.ActiveSessions
		}
	}
	if len(items) > 0 {
		n := float64(len(items))
		stats.AvgCPU /= n
		stats.AvgMemory /= n
		stats.AvgSessions /= n
		stats.AvgLatency /= n
	}
	if stats.TotalRequests > 0 {
		stats.SuccessRate = float64(stats.TotalSuccess) / float64(stats.TotalRequests) * 100
	}
	return &types.ProxyMetricsStatsResp{Msg: "success", Data: stats}, nil
}
