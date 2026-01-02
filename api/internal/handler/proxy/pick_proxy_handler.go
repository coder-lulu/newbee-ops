package proxy

import (
    "net/http"

    "github.com/coder-lulu/newbee-ops-api/internal/metrics"
    logic "github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
    "github.com/coder-lulu/newbee-ops-api/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// Pick 简单返回一个可用的 Proxy（临时占位）
func Pick(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query()
        proto := q.Get("protocol")
        region := q.Get("region")
        az := q.Get("az")
        plg := logic.NewPickLogic(svcCtx.ProxyRegistry, svcCtx.Config.Ops.DefaultProxyEndpoints)
        proxyId, endpoint, ok := plg.Pick(proto, region, az)
        if ok { metrics.IncPick() }
        httpx.OkJsonCtx(r.Context(), w, types.PickProxyResp{
            Code: 0,
            Msg:  "success",
            Data: types.PickProxyData{
                ProxyId:  proxyId,
                Endpoint: endpoint,
            },
        })
    }
}
