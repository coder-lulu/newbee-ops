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
	// todo: add your logic here and delete this line

	return
}
