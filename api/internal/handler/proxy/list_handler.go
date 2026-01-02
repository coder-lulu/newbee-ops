package proxy

import (
    "net/http"

    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// List 列出注册表内的 Proxy（简化输出）
func List(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lg := logic.NewListLogic(svcCtx.ProxyRegistry)
        httpx.OkJsonCtx(r.Context(), w, lg.List())
    }
}
