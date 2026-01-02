package proxy

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Proxy心跳（PSK认证，无需JWT）
func ProxyHeartbeatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProxyHeartbeatReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxy.NewProxyHeartbeatLogic(r.Context(), svcCtx)
		resp, err := l.ProxyHeartbeat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
