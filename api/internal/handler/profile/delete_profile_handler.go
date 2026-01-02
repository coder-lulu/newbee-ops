package profile

import (
    "net/http"

    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/profile"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func Delete(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ciId := r.URL.Query().Get("ciId")
        if ciId == "" { http.Error(w, "ciId required", http.StatusBadRequest); return }
        lg := logic.NewLogic(r.Context(), svcCtx)
        if err := lg.Delete(ciId); err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        httpx.OkJsonCtx(r.Context(), w, types.BaseMsgResp{Code: 0, Msg: "profile deleted successfully"})
    }
}
