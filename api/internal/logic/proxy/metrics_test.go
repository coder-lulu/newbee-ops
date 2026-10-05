package proxy

import (
	"context"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"testing"
)

type sampleClient struct {
	rpcclient.OpsClient
	seen context.Context
}

func (c *sampleClient) GetProxyMetricsList(ctx context.Context, r *pb.ProxyMetricsListReq) (*pb.ProxyMetricsListResp, error) {
	c.seen = ctx
	stamp := int64(100000)
	cpu := float64(40)
	req := uint64(5)
	ok := uint64(4)
	return &pb.ProxyMetricsListResp{Total: 1, Data: []*pb.ProxyMetricsInfo{{ProxyId: r.ProxyId, Timestamp: &stamp, CpuUsage: &cpu, RequestCountDelta: &req, SuccessCountDelta: &ok}}}, nil
}
func TestMetricsReadFiltersAndAggregates(t *testing.T) {
	client := &sampleClient{}
	ctx := context.WithValue(context.Background(), struct{}{}, "tenant")
	got, err := readMetrics(ctx, client, "demo", 101, 200)
	if err != nil || len(got) != 0 {
		t.Fatalf("time filter lost: %+v %v", got, err)
	}
	if client.seen != ctx {
		t.Fatal("tenant context lost")
	}
	result, err := NewGetProxyMetricsStatsLogic(ctx, &svc.ServiceContext{OpsClient: client}).GetProxyMetricsStats(&types.ProxyMetricsStatsReq{ProxyID: "demo", StartTime: 99, EndTime: 100})
	if err != nil || result.Data.DataPoints != 1 || result.Data.AvgCPU != 40 || result.Data.TotalRequests != 5 || result.Data.SuccessRate != 80 {
		t.Fatalf("wrong stats: %+v %v", result, err)
	}
	if _, err = readMetrics(ctx, client, "", 0, 0); err == nil {
		t.Fatal("empty proxy must not read all tenants/proxies")
	}
}
