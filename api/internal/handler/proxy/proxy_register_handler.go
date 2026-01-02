package proxy

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Proxy注册（PSK认证，无需JWT）
func ProxyRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProxyRegisterReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxy.NewProxyRegisterLogic(r.Context(), svcCtx)
		resp, err := l.ProxyRegister(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
