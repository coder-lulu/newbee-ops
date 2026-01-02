package agent

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/agent"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新Agent
func UpdateAgentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AgentInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewUpdateAgentLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAgent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
