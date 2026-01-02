package session

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/session"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// Get 返回指定会话详情（内存）
func Get(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ""
        // Try extract from path suffix (…/ops/session/{id})
        parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
        if len(parts) >= 3 && parts[len(parts)-2] == "session" {
            id = parts[len(parts)-1]
        }
        if id == "" {
            id = r.URL.Query().Get("id")
        }
        if id == "" { http.Error(w, "id required", http.StatusBadRequest); return }
        lg := logic.NewGetSessionLogic(r.Context(), svcCtx)
        sess, ok := lg.Get(id)
        if !ok { http.Error(w, "not found", http.StatusNotFound); return }
        httpx.OkJsonCtx(r.Context(), w, types.SessionItemResp{
            Code: 0,
            Msg:  "success",
            Data: types.SessionItem{
                Id:        strconv.FormatUint(sess.ID, 10),
                TenantId:  strconv.FormatUint(sess.TenantId, 10),
                UserId:    strconv.FormatUint(sess.UserId, 10),
                CiId:      sess.CiId,
                Protocol:  sess.Protocol,
                ProxyId:   sess.ProxyId,
                Endpoint:  sess.Endpoint,
                CreatedAt: sess.CreatedAt,
                ExpiresAt: sess.ExpiresAt,
                Status:    sess.Status,
                ClosedAt:  sess.ClosedAt,
            },
        })
    }
}
