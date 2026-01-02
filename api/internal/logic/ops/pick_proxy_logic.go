package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PickProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPickProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PickProxyLogic {
	return &PickProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PickProxyLogic) PickProxy(req *types.ProxyPickReq) (resp *types.ProxyPickResp, err error) {
	// 兼容性接口：重定向到 /proxy/proxy_pick
	l.Logger.Infow("Redirecting legacy /ops/pick_proxy to /proxy/proxy_pick",
		logx.Field("strategy", req.Strategy),
		logx.Field("capabilities", req.RequiredCapabilities))

	// 委托给新的proxy logic
	proxyLogic := proxy.NewProxyPickLogic(l.ctx, l.svcCtx)
	return proxyLogic.ProxyPick(req)
}
