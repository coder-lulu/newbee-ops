package profile

import (
    "net/http"

    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/profile"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func Get(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ciId := r.URL.Query().Get("ciId")
        if ciId == "" { http.Error(w, "ciId required", http.StatusBadRequest); return }
        lg := logic.NewLogic(r.Context(), svcCtx)
        ap, ok, err := lg.Get(ciId)
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        if !ok { http.Error(w, "not found", http.StatusNotFound); return }
        httpx.OkJsonCtx(r.Context(), w, types.AccessProfileResp{
            Code: 0,
            Msg:  "success",
            Data: ap,
        })
    }
}

func toIntMap(in map[string]int32) map[string]int {
    if in == nil { return nil }
    out := make(map[string]int, len(in))
    for k, v := range in { out[k] = int(v) }
    return out
}
