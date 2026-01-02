package ops

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/ops"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ops.NewListProfileLogic(r.Context(), svcCtx)
		resp, err := l.ListProfile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
