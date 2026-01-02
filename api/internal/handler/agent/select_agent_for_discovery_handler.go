package agent

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/agent"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 选择Agent（供Discovery使用）
func SelectAgentForDiscoveryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AgentSelectionReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewSelectAgentForDiscoveryLogic(r.Context(), svcCtx)
		resp, err := l.SelectAgentForDiscovery(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
