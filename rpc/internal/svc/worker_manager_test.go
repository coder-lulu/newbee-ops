package svc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
	_ "modernc.org/sqlite"
)

func TestWorkerManagerBootstrapThenRecovery(t *testing.T) {
	db, err := sql.Open("sqlite", "file:worker-bootstrap?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, db)))
	t.Cleanup(func() { _ = client.Close() })
	if err := hooks.QuickSetup(client); err != nil {
		t.Fatal(err)
	}
	manager, err := newWorkerManager(client, false)
	if err != nil || manager == nil {
		t.Fatalf("empty database bootstrap: manager=%v err=%v", manager, err)
	}
	if len(manager.ListWorkers(1)) != 0 {
		t.Fatal("bootstrap registry must be empty")
	}
	manager.Stop()
	if _, err := newWorkerManager(client, true); err == nil {
		t.Fatal("enabled startup must report missing tables")
	}
	ctx := context.Background()
	if err := client.Schema.Create(ctx, schema.WithForeignKeys(false)); err != nil {
		t.Fatal(err)
	}
	tenantCtx := hooks.SetTenantIDToContext(ctx, 1)
	_, err = client.Proxy.Create().SetWorkerID("bootstrap-worker").SetName("bootstrap-worker").SetIP("127.0.0.1").
		SetWorkerStatus(proxy.WorkerStatusOnline).SetLastHeartbeat(time.Now()).Save(tenantCtx)
	if err != nil {
		t.Fatal(err)
	}
	manager, err = newWorkerManager(client, true)
	if err != nil {
		t.Fatalf("startup after initialization: %v", err)
	}
	if len(manager.ListWorkers(1)) != 1 || len(manager.ListWorkers(2)) != 0 {
		t.Fatal("recovery lost worker or tenant isolation")
	}
	manager.Stop()
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := newWorkerManager(client, true); err == nil {
		t.Fatal("enabled startup must report database connection errors")
	}
}
