package script

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/script"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
)

// swagger:route post /script_version script GetScriptVersionById
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: IDReq
//
// Responses:
//  200: ScriptVersionDetailResp

func GetScriptVersionByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IDReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := script.NewGetScriptVersionByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetScriptVersionById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
