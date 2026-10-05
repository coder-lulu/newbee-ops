package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/coder-lulu/newbee-common/v2/middleware/keys"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-ops-api/internal/rpcclient"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type registrationRPC struct {
	rpcclient.OpsClient
	contexts []context.Context
	existing bool
	writes   []*ops.ProxyInfo
}

func TestRegistrationContextRPCPropagation(t *testing.T) {
	s := &svc.ServiceContext{OpsClient: &registrationRPC{}}
	s.Config.Ops.Registration.PSK = "test-psk"
	s.Config.Ops.Registration.TenantID = 23
	request, cancel := context.WithCancel(context.Background())
	bound, rejected := registrationContext(request, s, "test-psk")
	if rejected != nil {
		t.Fatal(rejected)
	}
	async := context.WithoutCancel(bound)
	cancel()
	if async.Err() != nil {
		t.Fatal("asynchronous context canceled with request")
	}
	for _, ctx := range []context.Context{bound, async} {
		err := hooks.SystemContextClientInterceptor()(ctx, "/ops.Proxy/Test", nil, nil, nil,
			func(out context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
				md, _ := metadata.FromOutgoingContext(out)
				_, err := hooks.ContextPropagationServerInterceptor()(metadata.NewIncomingContext(context.Background(), md), nil, &grpc.UnaryServerInfo{FullMethod: "/ops.Proxy/Test"},
					func(in context.Context, _ any) (any, error) {
						if hooks.GetCurrentTenantID(in) != 23 || hooks.IsSystemContext(in) {
							t.Error("RPC server lost tenant or bypassed isolation")
						}
						return nil, nil
					})
				return err
			})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (r *registrationRPC) GetProxyByProxyId(ctx context.Context, _ string) (*ops.ProxyInfo, error) {
	r.contexts = append(r.contexts, ctx)
	if r.existing {
		return &ops.ProxyInfo{}, nil
	}
	return nil, nil
}
func (r *registrationRPC) CreateProxy(ctx context.Context, info *ops.ProxyInfo) (*ops.BaseIDResp, error) {
	r.contexts = append(r.contexts, ctx)
	r.writes = append(r.writes, info)
	return &ops.BaseIDResp{}, nil
}
func (r *registrationRPC) UpdateProxy(ctx context.Context, info *ops.ProxyInfo) (*ops.BaseResp, error) {
	r.contexts = append(r.contexts, ctx)
	r.writes = append(r.writes, info)
	return &ops.BaseResp{}, nil
}

func TestRegistrationTenantBoundary(t *testing.T) {
	for _, action := range []string{"register", "heartbeat"} {
		for _, tc := range []struct {
			name, configuredPSK, requestPSK string
			tenant                          uint64
			want                            uint32
			nilRPC                          bool
		}{
			{"wrong PSK", "test-psk", "wrong", 17, 401, false},
			{"empty PSK", "", "", 17, 401, false},
			{"unbound tenant", "test-psk", "test-psk", 0, 503, false},
			{"missing RPC", "test-psk", "test-psk", 17, 503, true},
			{"bound tenant", "test-psk", "test-psk", 17, 0, false},
		} {
			t.Run(action+"/"+tc.name, func(t *testing.T) {
				fake := &registrationRPC{existing: action == "heartbeat"}
				s := &svc.ServiceContext{OpsClient: fake}
				if tc.nilRPC {
					s.OpsClient = nil
				}
				s.Config.Ops.Registration.PSK = tc.configuredPSK
				s.Config.Ops.Registration.TenantID = tc.tenant
				ctx := keys.NewContextManager().SetTenantID(context.Background(), "999")
				ctx = context.WithValue(ctx, keys.SystemContextKey, true)
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(keys.TenantIDKey.String(), "998", keys.SystemContextKey.String(), "true"))
				ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(keys.TenantIDKey.String(), "997", keys.SystemContextKey.String(), "true"))
				var resp *types.BaseResp
				var err error
				if action == "register" {
					resp, err = NewProxyRegisterLogic(ctx, s).ProxyRegister(&types.ProxyRegisterReq{ProxyID: "local", PSK: tc.requestPSK})
				} else {
					resp, err = NewProxyHeartbeatLogic(ctx, s).ProxyHeartbeat(&types.ProxyHeartbeatReq{ProxyID: "local", PSK: tc.requestPSK, ProxyStatus: "online"})
				}
				if err != nil || resp == nil || resp.Code != tc.want {
					t.Fatalf("response=%v err=%v want code=%d", resp, err, tc.want)
				}
				if tc.want != 0 {
					if len(fake.contexts) != 0 {
						t.Fatalf("rejected request called RPC %d times", len(fake.contexts))
					}
					return
				}
				if len(fake.contexts) != 2 {
					t.Fatalf("RPC calls=%d", len(fake.contexts))
				}
				for _, info := range fake.writes {
					if info.LastHeartbeat == nil || time.Since(time.UnixMilli(*info.LastHeartbeat)).Abs() > 5*time.Second {
						t.Error("heartbeat must use current Unix milliseconds")
					}
					if action == "register" && (info.RegisterTime == nil || time.Since(time.UnixMilli(*info.RegisterTime)).Abs() > 5*time.Second) {
						t.Error("registration must use current Unix milliseconds")
					}
				}
				for _, got := range fake.contexts {
					if keys.NewContextManager().GetTenantID(got) != "17" || hooks.IsSystemContext(got) {
						t.Error("trusted tenant binding or isolation missing")
					}
					md, _ := metadata.FromOutgoingContext(got)
					if values := md.Get(keys.TenantIDKey.String()); len(values) != 1 || values[0] != "17" {
						t.Errorf("outgoing tenant=%v", values)
					}
					if len(md.Get(keys.SystemContextKey.String())) != 0 {
						t.Error("outgoing bypass flag")
					}
				}
			})
		}
	}
}
