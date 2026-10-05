package proxy

import (
	"context"
	"errors"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"testing"
)

type readProxyClient struct {
	rpcclient.OpsClient
	request  *pb.ProxyListReq
	ctx      context.Context
	response *pb.ProxyListResp
	err      error
}

func (c *readProxyClient) GetProxyList(ctx context.Context, r *pb.ProxyListReq) (*pb.ProxyListResp, error) {
	c.ctx = ctx
	c.request = r
	return c.response, c.err
}
func TestGetProxyListPreservesContextAndFilters(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "tenant-context")
	id := uint64(7)
	name := "演示代理"
	state := "offline"
	client := &readProxyClient{response: &pb.ProxyListResp{Total: 31, Data: []*pb.ProxyInfo{{Id: &id, Name: &name, ProxyStatus: &state}, {}}}}
	result, err := NewGetProxyListLogic(ctx, &svc.ServiceContext{OpsClient: client}).GetProxyList(&types.ProxyListReq{Page: 2, PageSize: 20, Status: 2, ProxyStatus: "offline", Name: "演示", Region: "cn-beijing", Zone: "demo", Tags: "demo"})
	if err != nil {
		t.Fatal(err)
	}
	if client.ctx != ctx || client.request.Page != 2 || client.request.PageSize != 20 || client.request.GetStatus() != 2 || client.request.GetProxyStatus() != "offline" || client.request.GetRegion() != "cn-beijing" || client.request.GetZone() != "demo" || client.request.GetTags() != "demo" {
		t.Fatalf("request/context lost: %+v", client.request)
	}
	if result.Data.Total != 31 || len(result.Data.Data) != 2 || result.Data.Data[0].Name != name || result.Data.Data[0].Id != id || result.Data.Data[0].ProxyStatus != state {
		t.Fatalf("wrong response: %+v", result)
	}
}
func TestGetProxyListErrorAndEmpty(t *testing.T) {
	expected := errors.New("denied")
	client := &readProxyClient{err: expected}
	logic := NewGetProxyListLogic(context.Background(), &svc.ServiceContext{OpsClient: client})
	if _, err := logic.GetProxyList(&types.ProxyListReq{}); !errors.Is(err, expected) {
		t.Fatalf("RPC error swallowed: %v", err)
	}
	client.err = nil
	client.response = &pb.ProxyListResp{}
	result, err := logic.GetProxyList(&types.ProxyListReq{})
	if err != nil || result.Data.Data == nil {
		t.Fatalf("empty list must be [], got %+v %v", result, err)
	}
	if client.request.Status != nil || client.request.Name != nil || client.request.ProxyStatus != nil {
		t.Fatal("absent filters became constraints")
	}
}
