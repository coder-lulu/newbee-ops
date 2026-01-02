package proxygroup

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxygroup"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建Proxy分组（需要JWT）
func CreateProxyGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProxyGroupInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxygroup.NewCreateProxyGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateProxyGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
