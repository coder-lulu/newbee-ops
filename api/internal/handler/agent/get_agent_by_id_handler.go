package agent

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/agent"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取Agent详情
func GetAgentByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IDPathReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewGetAgentByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetAgentById(&types.IDReq{Id: req.Id})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, map[string]any{"code": 0, "msg": "success", "data": resp})
		}
	}
}
