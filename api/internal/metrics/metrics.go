package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    proxyRegistered = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "ops_center",
        Subsystem: "proxy",
        Name:      "registered_total",
        Help:      "Number of registered proxies (TTL valid)",
    })
    proxyHeartbeatTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "ops_center",
        Subsystem: "proxy",
        Name:      "heartbeat_total",
        Help:      "Total heartbeats received",
    })
    proxyPickTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "ops_center",
        Subsystem: "proxy",
        Name:      "pick_total",
        Help:      "Total picks made for sessions",
    })
)

func init() {
    prometheus.MustRegister(proxyRegistered, proxyHeartbeatTotal, proxyPickTotal)
}

func SetRegisteredGauge(n int) { proxyRegistered.Set(float64(n)) }
func IncHeartbeat()           { proxyHeartbeatTotal.Inc() }
func IncPick()                { proxyPickTotal.Inc() }

