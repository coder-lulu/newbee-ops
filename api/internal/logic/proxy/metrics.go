package proxy

import (
	"context"
	"fmt"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"sort"
)

// Read persisted samples through the tenant-aware RPC client; never query a proxy host.
func readMetrics(ctx context.Context, client rpcclient.OpsClient, proxyID string, start, end int64) ([]types.ProxyMetricsItem, error) {
	if proxyID == "" {
		return nil, fmt.Errorf("proxy_id is required")
	}
	if start > 0 && end > 0 && start > end {
		return nil, fmt.Errorf("start_time must not exceed end_time")
	}
	items := make([]types.ProxyMetricsItem, 0)
	for page := uint64(1); ; page++ {
		result, err := client.GetProxyMetricsList(ctx, &pb.ProxyMetricsListReq{Page: page, PageSize: 1000, ProxyId: &proxyID})
		if err != nil {
			return nil, err
		}
		for _, x := range result.Data {
			stamp := x.GetTimestamp() / 1000
			if start > 0 && stamp < start || end > 0 && stamp > end {
				continue
			}
			items = append(items, types.ProxyMetricsItem{ProxyID: x.GetProxyId(), Timestamp: stamp, CPUUsage: x.GetCpuUsage(), MemoryUsage: x.GetMemoryUsage(), DiskUsage: x.GetDiskUsage(), NetworkInDelta: x.GetNetworkInDelta(), NetworkOutDelta: x.GetNetworkOutDelta(), ActiveSessions: int(x.GetActiveSessions()), RequestDelta: x.GetRequestCountDelta(), SuccessDelta: x.GetSuccessCountDelta(), FailureDelta: x.GetFailureCountDelta(), AvgLatencyMs: x.GetAvgLatencyMs(), ProxyStatus: x.GetProxyStatus()})
		}
		if len(result.Data) == 0 || page*1000 >= result.Total {
			break
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Timestamp > items[j].Timestamp })
	return items, nil
}
