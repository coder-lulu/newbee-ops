package session

import (
    "net/http"

    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/session"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// Create 创建会话并签发短期Session Token
func Create(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.CreateSessionReq
        if err := httpx.Parse(r, &req, true); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }
        lg := logic.NewCreateSessionLogic(r.Context(), svcCtx)
        resp, err := lg.Create(&req)
        if err != nil {
            if logic.IsBadRequest(err) { http.Error(w, err.Error(), http.StatusBadRequest); return }
            httpx.ErrorCtx(r.Context(), w, err); return
        }
        httpx.OkJsonCtx(r.Context(), w, resp)
    }
}
