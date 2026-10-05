package ops

import (
	"context"
	"encoding/json"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"strings"
	"testing"
)

type profileReadClient struct {
	rpcclient.OpsClient
	ctx     context.Context
	request *pb.AccessProfileListReq
}

func (c *profileReadClient) GetAccessProfileList(ctx context.Context, r *pb.AccessProfileListReq) (*pb.AccessProfileListResp, error) {
	c.ctx = ctx
	c.request = r
	id := "demo-ci"
	return &pb.AccessProfileListResp{Total: 1, Data: []*pb.AccessProfileInfo{{CiId: &id}}}, nil
}
func TestProfileReadsMatchUIContract(t *testing.T) {
	client := &profileReadClient{}
	ctx := context.WithValue(context.Background(), struct{}{}, "tenant")
	svcCtx := &svc.ServiceContext{OpsClient: client}
	list, err := NewListProfileLogic(ctx, svcCtx).ListProfile()
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"items":[{"ciId":"demo-ci"`) {
		t.Fatalf("UI cannot read profile fields: %s", body)
	}
	detail, err := NewGetProfileLogic(ctx, svcCtx).GetProfile(&types.AccessProfileQueryReq{CiId: "demo-ci"})
	if err != nil || detail.Data.CiId != "demo-ci" || client.request.GetCiId() != "demo-ci" || client.ctx != ctx {
		t.Fatalf("detail query/context lost: %+v %v", detail, err)
	}
	if _, err = NewGetProfileLogic(ctx, svcCtx).GetProfile(&types.AccessProfileQueryReq{}); err == nil {
		t.Fatal("missing CI must not return an arbitrary profile")
	}
}
