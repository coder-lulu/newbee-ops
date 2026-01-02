package proxy_metrics

import (
	"context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/proxymetrics"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteProxyMetricsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProxyMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProxyMetricsLogic {
	return &DeleteProxyMetricsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProxyMetricsLogic) DeleteProxyMetrics(in *ops.IDsReq) (*ops.BaseResp, error) {
	_, err := l.svcCtx.DB.ProxyMetrics.Delete().Where(proxymetrics.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.DeleteSuccess }, nil
}
