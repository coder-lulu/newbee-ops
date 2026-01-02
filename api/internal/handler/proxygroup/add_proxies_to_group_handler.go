package proxygroup

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxygroup"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 添加Proxy到分组（需要JWT）
func AddProxiesToGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddProxiesToGroupReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxygroup.NewAddProxiesToGroupLogic(r.Context(), svcCtx)
		resp, err := l.AddProxiesToGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
