package handler

import (
	"context"
	"encoding/json"
	"errors"
	ah "github.com/coder-lulu/newbee-ops-api/internal/handler/agent"
	agh "github.com/coder-lulu/newbee-ops-api/internal/handler/agentgroup"
	oh "github.com/coder-lulu/newbee-ops-api/internal/handler/ops"
	ph "github.com/coder-lulu/newbee-ops-api/internal/handler/proxy"
	pgh "github.com/coder-lulu/newbee-ops-api/internal/handler/proxygroup"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"net/http"
	"net/http/httptest"
	"testing"
)

type readContractClient struct {
	rpcclient.OpsClient
	id  uint64
	err error
}

func (c *readContractClient) GetAgentList(context.Context, *pb.AgentListReq) (*pb.AgentListResp, error) {
	return &pb.AgentListResp{Total: 30, Data: []*pb.AgentInfo{{}}}, c.err
}
func (c *readContractClient) GetAgentGroupList(context.Context, *pb.AgentGroupListReq) (*pb.AgentGroupListResp, error) {
	return &pb.AgentGroupListResp{Total: 30, Data: []*pb.AgentGroupInfo{{}}}, c.err
}
func (c *readContractClient) GetProxyById(_ context.Context, r *pb.IDReq) (*pb.ProxyInfo, error) {
	c.id = r.Id
	return &pb.ProxyInfo{Id: &r.Id}, c.err
}
func (c *readContractClient) GetAgentById(_ context.Context, r *pb.IDReq) (*pb.AgentInfo, error) {
	c.id = r.Id
	return &pb.AgentInfo{Id: &r.Id}, c.err
}
func (c *readContractClient) GetAgentGroupById(_ context.Context, r *pb.IDReq) (*pb.AgentGroupInfo, error) {
	c.id = r.Id
	return &pb.AgentGroupInfo{Id: &r.Id}, c.err
}
func (c *readContractClient) GetProxyGroupById(_ context.Context, r *pb.IDReq) (*pb.ProxyGroupInfo, error) {
	c.id = r.Id
	return &pb.ProxyGroupInfo{Id: &r.Id}, c.err
}
func TestReadHandlersMatchFrontendEnvelope(t *testing.T) {
	for _, tc := range []struct {
		name   string
		create func(*svc.ServiceContext) http.HandlerFunc
		detail bool
	}{
		{"agent-list", ah.GetAgentListHandler, false}, {"agent-group-list", agh.GetAgentGroupListHandler, false},
		{"proxy-detail", ph.GetProxyByIdHandler, true}, {"agent-detail", ah.GetAgentByIdHandler, true}, {"agent-group-detail", agh.GetAgentGroupByIdHandler, true}, {"proxy-group-detail", pgh.GetProxyGroupByIdHandler, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &readContractClient{}
			h := tc.create(&svc.ServiceContext{OpsClient: client})
			request := httptest.NewRequest(http.MethodGet, "/read/73?page=1&pageSize=20", nil)
			if tc.detail {
				request = pathvar.WithVars(request, map[string]string{"id": "73"})
			}
			recorder := httptest.NewRecorder()
			h(recorder, request)
			var envelope struct {
				Code *int                       `json:"code"`
				Data map[string]json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
				t.Fatal(err)
			}
			// requestClient accepts only code===0 and returns envelope.data exactly once.
			if recorder.Code != 200 || envelope.Code == nil || *envelope.Code != 0 {
				t.Fatalf("frontend would reject: %s", recorder.Body.String())
			}
			if tc.detail {
				if client.id != 73 || string(envelope.Data["id"]) != "73" {
					t.Fatalf("path ID or detail unwrap lost: %s", recorder.Body.String())
				}
			} else {
				if string(envelope.Data["total"]) != "30" || len(envelope.Data["data"]) == 0 {
					t.Fatalf("list requires one unwrap: %s", recorder.Body.String())
				}
			}
			client.err = errors.New("RPC unavailable")
			recorder = httptest.NewRecorder()
			h(recorder, request)
			if recorder.Code == 200 {
				t.Fatalf("RPC failure became success: %s", recorder.Body.String())
			}
		})
	}
}

func (c *readContractClient) GetSessionList(_ context.Context, r *pb.SessionListReq) (*pb.SessionListResp, error) {
	c.id = r.PageSize
	sid := "demo-session-01"
	return &pb.SessionListResp{Total: 30, Data: []*pb.SessionInfo{{SessionId: &sid}}}, c.err
}
func TestSessionListAcceptsOptionalQueryFilters(t *testing.T) {
	client := &readContractClient{}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ops/session/list?page=1&size=20", nil)
	oh.ListSessionHandler(&svc.ServiceContext{OpsClient: client})(recorder, request)
	var response struct {
		Code *int `json:"code"`
		Data struct {
			Items []map[string]any `json:"items"`
			Total int              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid response: %s", recorder.Body.String())
	}
	if response.Code == nil || *response.Code != 0 || client.id != 20 || response.Data.Total != 30 || len(response.Data.Items) != 1 {
		t.Fatalf("optional session query rejected: %s", recorder.Body.String())
	}
	if response.Data.Items[0]["id"] != "demo-session-01" {
		t.Fatal("real session identifier lost")
	}
}
