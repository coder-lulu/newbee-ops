package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterProxyLogic {
	return &RegisterProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterProxyLogic) RegisterProxy(req *types.ProxyRegisterReq) (resp *types.BaseResp, err error) {
	// 兼容性接口：重定向到 /proxy/proxy_register
	l.Logger.Infow("Redirecting legacy /ops/register_proxy to /proxy/proxy_register",
		logx.Field("proxy_id", req.ProxyID))

	// 委托给新的proxy logic
	proxyLogic := proxy.NewProxyRegisterLogic(l.ctx, l.svcCtx)
	return proxyLogic.ProxyRegister(req)
}
