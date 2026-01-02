package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHeartbeatProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatProxyLogic {
	return &HeartbeatProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeartbeatProxyLogic) HeartbeatProxy(req *types.ProxyHeartbeatReq) (resp *types.BaseResp, err error) {
	// 兼容性接口：重定向到 /proxy/proxy_heartbeat
	l.Logger.Infow("Redirecting legacy /ops/heartbeat_proxy to /proxy/proxy_heartbeat",
		logx.Field("proxy_id", req.ProxyID))

	// 委托给新的proxy logic
	proxyLogic := proxy.NewProxyHeartbeatLogic(l.ctx, l.svcCtx)
	return proxyLogic.ProxyHeartbeat(req)
}
