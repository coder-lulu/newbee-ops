package agent

import (
	"context"
	entcore "entgo.io/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"testing"
	"time"
)

func TestOfflineAgentTimeMapping(t *testing.T) {
	now := time.UnixMilli(1791200000123)
	zero := time.Time{}
	for _, stamp := range []*time.Time{nil, &zero, &now} {
		for _, route := range []string{"list", "detail"} {
			client := ent.NewClient()
			row := &ent.Agent{ID: 7, LastHeartbeat: stamp, LastOnlineAt: stamp}
			client.Agent.Intercept(entcore.InterceptFunc(func(entcore.Querier) entcore.Querier {
				return entcore.QuerierFunc(func(ctx context.Context, _ entcore.Query) (entcore.Value, error) {
					if entcore.QueryFromContext(ctx).Op == entcore.OpQueryCount {
						return 1, nil
					}
					return []*ent.Agent{row}, nil
				})
			}))
			service := &svc.ServiceContext{DB: client}
			var got *ops.AgentInfo
			if route == "list" {
				res, err := NewGetAgentListLogic(context.Background(), service).GetAgentList(&ops.AgentListReq{Page: 1, PageSize: 10})
				if err != nil {
					t.Fatal(err)
				}
				got = res.Data[0]
			} else {
				var err error
				got, err = NewGetAgentByIdLogic(context.Background(), service).GetAgentById(&ops.IDReq{Id: 7})
				if err != nil {
					t.Fatal(err)
				}
			}
			for _, value := range []*int64{got.LastHeartbeat, got.LastOnlineAt} {
				if stamp == nil || stamp.IsZero() {
					if value != nil {
						t.Fatal("absent timestamp became a value")
					}
				} else if value == nil || *value != stamp.UnixMilli() {
					t.Fatal("milliseconds lost")
				}
			}
		}
	}
}
