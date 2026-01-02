package session

import (
    "net/http"

    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/session"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// Close 关闭会话（内存存储）
func Close(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.CloseSessionReq
        if err := httpx.Parse(r, &req, true); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }
        if req.SessionId == "" { http.Error(w, "sessionId required", http.StatusBadRequest); return }
        lg := logic.NewCloseSessionLogic(r.Context(), svcCtx)
        ok, closedAt, msg := lg.Close(req.SessionId)
        if !ok {
            // 如果未找到或租户不匹配，返回标准错误响应
            httpx.OkJsonCtx(r.Context(), w, types.CloseSessionResp{
                Code: 1,
                Msg:  msg,
                Data: types.CloseSessionData{Ok: false, ClosedAt: 0},
            })
            return
        }
        httpx.OkJsonCtx(r.Context(), w, types.CloseSessionResp{
            Code: 0,
            Msg:  "session closed successfully",
            Data: types.CloseSessionData{Ok: true, ClosedAt: closedAt},
        })
    }
}
