package proxygroup

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxygroup"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 从分组选择Proxy（需要JWT）
func PickProxyFromGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PickProxyFromGroupReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxygroup.NewPickProxyFromGroupLogic(r.Context(), svcCtx)
		resp, err := l.PickProxyFromGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
