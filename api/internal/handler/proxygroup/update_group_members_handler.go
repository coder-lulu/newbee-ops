package proxygroup

import (
	"net/http"

	"github.com/coder-lulu/newbee-ops-api/internal/logic/proxygroup"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量更新分组成员（需要JWT）- 替换式更新，自动处理增删
func UpdateGroupMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGroupMembersReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proxygroup.NewUpdateGroupMembersLogic(r.Context(), svcCtx)
		resp, err := l.UpdateGroupMembers(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
