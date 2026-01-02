package session

import (
    "net/http"
    "strconv"

    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/session"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// List 返回会话列表，支持按 status/ciId 过滤与分页
// GET /ops/session/list?status=active|closed&ciId=&page=&size=
func List(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query()
        status := q.Get("status")
        ciId := q.Get("ciId")
        page, _ := strconv.Atoi(q.Get("page"))
        size, _ := strconv.Atoi(q.Get("size"))
        lg := logic.NewListSessionLogic(r.Context(), svcCtx)
        items, total := lg.List(status, ciId, page, size)
        // shape for response
        out := make([]types.SessionItem, 0, len(items))
        for _, s := range items {
            out = append(out, types.SessionItem{
                Id:        strconv.FormatUint(s.ID, 10),
                TenantId:  strconv.FormatUint(s.TenantId, 10),
                UserId:    strconv.FormatUint(s.UserId, 10),
                CiId:      s.CiId,
                Protocol:  s.Protocol,
                ProxyId:   s.ProxyId,
                Endpoint:  s.Endpoint,
                CreatedAt: s.CreatedAt,
                ExpiresAt: s.ExpiresAt,
                Status:    s.Status,
                ClosedAt:  s.ClosedAt,
            })
        }
        httpx.OkJsonCtx(r.Context(), w, types.SessionListResp{
            Code: 0,
            Msg:  "success",
            Data: types.SessionListData{Items: out, Total: total},
        })
    }
}
