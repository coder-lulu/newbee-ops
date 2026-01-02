package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyMetricsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询Proxy指标（需要JWT）
func NewGetProxyMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyMetricsLogic {
	return &GetProxyMetricsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyMetricsLogic) GetProxyMetrics(req *types.ProxyMetricsReq) (resp *types.ProxyMetricsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
