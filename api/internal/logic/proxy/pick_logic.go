package proxy

import proxyservice "github.com/coder-lulu/newbee-ops-api/internal/services/proxy"

type PickLogic struct {
    r                *proxyservice.Registry
    defaultEndpoints []string
}
func NewPickLogic(r *proxyservice.Registry, defaults []string) *PickLogic { return &PickLogic{r: r, defaultEndpoints: defaults} }

func (l *PickLogic) Pick(protocol, region, az string) (proxyId, endpoint string, ok bool) {
    if l.r != nil {
        var picked *proxyservice.ProxyInfo
        if protocol != "" || region != "" || az != "" {
            picked = l.r.PickWithFilter(protocol, region, az)
        } else {
            picked = l.r.Pick()
        }
        if picked != nil {
            ep := picked.Endpoints["http"]
            if ep == "" { ep = picked.Endpoints["ws"] }
            return picked.ProxyID, ep, true
        }
    }
    endpoint = "http://127.0.0.1:8889"
    proxyId = "proxy-default"
    if len(l.defaultEndpoints) > 0 { endpoint = l.defaultEndpoints[0]; proxyId = "proxy-001" }
    return proxyId, endpoint, false
}
