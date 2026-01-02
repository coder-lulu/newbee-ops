package ops

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/ops"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListSessionReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ops.NewListSessionLogic(r.Context(), svcCtx)
		resp, err := l.ListSession(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
