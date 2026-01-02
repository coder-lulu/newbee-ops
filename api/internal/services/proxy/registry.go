package proxy

import (
    "sync"
    "time"
)

// ProxyInfo 保存注册到中心的 Proxy 节点信息
type ProxyInfo struct {
    ProxyID    string            `json:"proxyId"`
    Region     string            `json:"region"`
    AZ         string            `json:"az"`
    Endpoints  map[string]string `json:"endpoints"` // e.g. {"ws":"ws://...","http":"http://..."}
    Capabilities []string        `json:"capabilities"`
    Labels     map[string]string `json:"labels"`
    Status     string            `json:"status"` // up|degraded|down
    Load       ProxyLoad         `json:"load"`
    Capacity   ProxyCapacity     `json:"capacity"`
    LastHeartbeat time.Time      `json:"lastHeartbeat"`
}

type ProxyLoad struct {
    CPU     float64 `json:"cpu"`
    Mem     float64 `json:"mem"`
    Session struct {
        Active int `json:"active"`
        Max    int `json:"max"`
    } `json:"session"`
}

type ProxyCapacity struct {
    Sessions int `json:"sessions"`
}

// Registry 简易内存注册表
type Registry struct {
    mu      sync.RWMutex
    items   map[string]*ProxyInfo
    ttl     time.Duration
}

func NewRegistry(ttl time.Duration) *Registry {
    if ttl <= 0 {
        ttl = 2 * time.Minute
    }
    return &Registry{items: make(map[string]*ProxyInfo), ttl: ttl}
}

func (r *Registry) Register(info *ProxyInfo) {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now()
    info.LastHeartbeat = now
    if info.Status == "" {
        info.Status = "up"
    }
    r.items[info.ProxyID] = info
}

func (r *Registry) Heartbeat(id string, status string, load ProxyLoad) (*ProxyInfo, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    it, ok := r.items[id]
    if !ok {
        return nil, false
    }
    it.Status = status
    it.Load = load
    it.LastHeartbeat = time.Now()
    return it, true
}

func (r *Registry) List() []*ProxyInfo {
    r.mu.RLock()
    defer r.mu.RUnlock()
    out := make([]*ProxyInfo, 0, len(r.items))
    cutoff := time.Now().Add(-r.ttl)
    for _, v := range r.items {
        if v.LastHeartbeat.Before(cutoff) {
            // 过期，跳过
            continue
        }
        out = append(out, v)
    }
    return out
}

// Pick 简易选择：优先 status=up，按会话活跃数升序
func (r *Registry) Pick() *ProxyInfo {
    lst := r.List()
    var best *ProxyInfo
    for _, it := range lst {
        if it.Status != "up" {
            continue
        }
        if best == nil || it.Load.Session.Active < best.Load.Session.Active {
            best = it
        }
    }
    if best != nil {
        return best
    }
    if len(lst) > 0 {
        return lst[0]
    }
    return nil
}

// PickWithFilter 带过滤条件选择
func (r *Registry) PickWithFilter(protocol, region, az string) *ProxyInfo {
    lst := r.List()
    var best *ProxyInfo
    for _, it := range lst {
        if it.Status != "up" { continue }
        if region != "" && it.Region != region { continue }
        if az != "" && it.AZ != az { continue }
        if protocol != "" {
            ok := false
            for _, c := range it.Capabilities { if c == protocol { ok = true; break } }
            if !ok { continue }
        }
        if best == nil || it.Load.Session.Active < best.Load.Session.Active { best = it }
    }
    if best != nil { return best }
    return r.Pick()
}
