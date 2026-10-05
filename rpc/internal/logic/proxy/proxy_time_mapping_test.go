package proxy

import (
	"context"
	"testing"
	"time"

	entcore "entgo.io/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

func TestProxyTimeMapping(t *testing.T) {
	now := time.UnixMilli(1791200000123)
	zero := time.Time{}
	for _, tc := range []struct {
		name                          string
		heartbeat, registered, health *time.Time
	}{
		{"new proxy without health check", &now, &now, nil},
		{"all null", nil, nil, nil},
		{"missing heartbeat", nil, &now, &now},
		{"missing registration", &now, nil, &now},
		{"all populated", &now, &now, &now},
		{"zero time", &zero, &zero, &zero},
	} {
		for _, route := range []string{"list", "detail"} {
			t.Run(tc.name+"/"+route, func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("mapping panicked: %v", p)
					}
				}()
				row := &ent.Proxy{ID: 7, WorkerID: "local", LastHeartbeat: tc.heartbeat, RegisterTime: tc.registered, LastHealthCheck: tc.health}
				client := ent.NewClient()
				// Stub query results at Ent's boundary; execute the real RPC mapping without a database.
				client.Proxy.Intercept(entcore.InterceptFunc(func(entcore.Querier) entcore.Querier {
					return entcore.QuerierFunc(func(ctx context.Context, _ entcore.Query) (entcore.Value, error) {
						if entcore.QueryFromContext(ctx).Op == entcore.OpQueryCount {
							return 1, nil
						}
						return []*ent.Proxy{row}, nil
					})
				}))
				s := &svc.ServiceContext{DB: client}
				var got *ops.ProxyInfo
				if route == "list" {
					resp, err := NewGetProxyListLogic(context.Background(), s).GetProxyList(&ops.ProxyListReq{Page: 1, PageSize: 10})
					if err != nil {
						t.Fatal(err)
					}
					if resp.Total != 1 || len(resp.Data) != 1 {
						t.Fatalf("unexpected list: %v", resp)
					}
					got = resp.Data[0]
				} else {
					var err error
					got, err = NewGetProxyByIdLogic(context.Background(), s).GetProxyById(&ops.IDReq{Id: 7})
					if err != nil {
						t.Fatal(err)
					}
				}
				for _, pair := range []struct {
					name  string
					src   *time.Time
					value *int64
				}{
					{"heartbeat", tc.heartbeat, got.LastHeartbeat}, {"registration", tc.registered, got.RegisterTime}, {"health", tc.health, got.LastHealthCheck},
				} {
					if pair.src == nil || pair.src.IsZero() {
						if pair.value != nil {
							t.Errorf("%s expected nil", pair.name)
						}
						continue
					}
					if pair.value == nil || *pair.value != pair.src.UnixMilli() {
						t.Errorf("%s lost millisecond value: %v", pair.name, pair.value)
					}
				}
			})
		}
	}
}
