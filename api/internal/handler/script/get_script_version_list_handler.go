package script

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/script"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
)

// swagger:route post /script_version/list script GetScriptVersionList
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ScriptVersionListReq
//
// Responses:
//  200: ScriptVersionListResp

func GetScriptVersionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ScriptVersionListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := script.NewGetScriptVersionListLogic(r.Context(), svcCtx)
		resp, err := l.GetScriptVersionList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
