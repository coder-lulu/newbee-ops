package profile

import (
    "net/http"

    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/profile"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func List(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lg := logic.NewLogic(r.Context(), svcCtx)
        items, total, err := lg.List()
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        httpx.OkJsonCtx(r.Context(), w, types.AccessProfileListResp{
            Code: 0,
            Msg:  "success",
            Data: types.AccessProfileListData{Items: items, Total: total},
        })
    }
}
