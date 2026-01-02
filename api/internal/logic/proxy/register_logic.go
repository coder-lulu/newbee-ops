package proxy

import (
	"github.com/coder-lulu/newbee-ops-api/internal/metrics"
	proxyservice "github.com/coder-lulu/newbee-ops-api/internal/services/proxy"
)

type RegisterLogic struct { r *proxyservice.Registry }
func NewRegisterLogic(r *proxyservice.Registry) *RegisterLogic { return &RegisterLogic{r: r} }
func (l *RegisterLogic) Register(info *proxyservice.ProxyInfo) bool {
    if l.r == nil { return false }
    l.r.Register(info)
    if size := len(l.r.List()); size >= 0 { metrics.SetRegisteredGauge(size) }
    return true
}

type HeartbeatLogic struct { r *proxyservice.Registry }
func NewHeartbeatLogic(r *proxyservice.Registry) *HeartbeatLogic { return &HeartbeatLogic{r: r} }
func (l *HeartbeatLogic) Heartbeat(id string, status string, load proxyservice.ProxyLoad) (accepted bool, reason string) {
    if l.r == nil { return false, "no registry" }
    if _, ok := l.r.Heartbeat(id, status, load); !ok { return false, "not registered" }
    metrics.IncHeartbeat(); return true, ""
}
