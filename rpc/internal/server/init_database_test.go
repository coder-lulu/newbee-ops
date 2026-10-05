package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/coder-lulu/newbee-common/v2/enum/common"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/coder-lulu/newbee-core/rpc/coreclient"
	"github.com/coder-lulu/newbee-core/rpc/types/core"
	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/migrate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

type bootstrapCore struct {
	coreclient.Core
	menus  []*core.MenuInfo
	grants []uint64
	apis   map[string]bool
	fail   bool
}

func (c *bootstrapCore) GetRoleList(ctx context.Context, in *core.RoleListReq, _ ...grpc.CallOption) (*core.RoleListResp, error) {
	if c.fail {
		return nil, fmt.Errorf("Core unavailable")
	}
	if hooks.GetCurrentTenantID(ctx) != 1 {
		return nil, fmt.Errorf("missing bootstrap tenant")
	}
	return &core.RoleListResp{Data: []*core.RoleInfo{{Id: pointy.GetPointer(uint64(9)), Code: pointy.GetPointer("superadmin")}}}, nil
}
func (c *bootstrapCore) GetMenuList(context.Context, *core.PageInfoReq, ...grpc.CallOption) (*core.MenuInfoList, error) {
	return &core.MenuInfoList{Data: c.menus}, nil
}
func (c *bootstrapCore) CreateMenu(_ context.Context, in *core.MenuInfo, _ ...grpc.CallOption) (*core.BaseIDResp, error) {
	if in.GetParentId() != common.DefaultParentId {
		found := false
		for _, menu := range c.menus {
			if menu.GetId() == in.GetParentId() {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("parent menu does not exist: %d", in.GetParentId())
		}
	}
	for _, menu := range c.menus {
		if menu.GetPath() == in.GetPath() {
			return nil, fmt.Errorf("duplicate menu path: %s", in.GetPath())
		}
		if menu.GetName() == in.GetName() && menu.GetMenuType() == in.GetMenuType() {
			return nil, fmt.Errorf("duplicate menu name: %s", in.GetName())
		}
	}
	id := uint64(1001 + len(c.menus))
	in.Id = &id
	c.menus = append(c.menus, in)
	return &core.BaseIDResp{Id: id}, nil
}
func (c *bootstrapCore) GetMenuAuthority(context.Context, *core.IDReq, ...grpc.CallOption) (*core.RoleMenuAuthorityResp, error) {
	return &core.RoleMenuAuthorityResp{MenuIds: c.grants}, nil
}
func (c *bootstrapCore) CreateOrUpdateMenuAuthority(_ context.Context, in *core.RoleMenuAuthorityReq, _ ...grpc.CallOption) (*core.BaseResp, error) {
	if in.RoleId != 9 {
		return nil, fmt.Errorf("unexpected role")
	}
	c.grants = in.MenuIds
	return &core.BaseResp{}, nil
}
func (c *bootstrapCore) CreateApi(_ context.Context, in *core.ApiInfo, _ ...grpc.CallOption) (*core.BaseIDResp, error) {
	c.apis[in.GetMethod()+" "+in.GetPath()] = true
	return &core.BaseIDResp{Id: 1}, nil
}

// Exercise the same generated server entry point used by the bootstrap RPC.
func TestInitDatabaseCreatesSchemaAndPreservesData(t *testing.T) {
	db, err := sql.Open("sqlite", "file:ops-bootstrap?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, db)))
	t.Cleanup(func() { _ = client.Close() })
	if err := hooks.QuickSetup(client); err != nil {
		t.Fatal(err)
	}
	coreRPC := &bootstrapCore{grants: []uint64{42}, apis: make(map[string]bool)}
	server := NewOpsServer(&svc.ServiceContext{DB: client, CoreRpc: coreRPC})
	ctx := context.Background()
	resp, err := server.InitDatabase(ctx, &ops.Empty{})
	if err != nil || resp.GetMsg() == "" {
		t.Fatalf("initialize: response=%v err=%v", resp, err)
	}
	for _, table := range migrate.Tables {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table.Name).Scan(&count); err != nil {
			t.Fatalf("table %s cannot be queried: %v", table.Name, err)
		}
	}
	tenantCtx := hooks.SetTenantIDToContext(ctx, 1)
	profile, err := client.AccessProfile.Create().SetCiID("bootstrap-ci").Save(tenantCtx)
	if err != nil {
		t.Fatal(err)
	}
	// Model another module sharing this database and a legacy Ops column.
	for _, ddl := range []string{
		"CREATE TABLE other_module (id INTEGER PRIMARY KEY)",
		"INSERT INTO other_module VALUES (7)",
		"ALTER TABLE access_profiles ADD COLUMN legacy_note TEXT",
		"CREATE INDEX legacy_profile_note ON access_profiles (legacy_note)",
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := server.InitDatabase(ctx, &ops.Empty{}); err != nil {
		t.Fatalf("repeat initialize: %v", err)
	}
	if len(coreRPC.menus) != 12 || len(coreRPC.grants) != 13 || coreRPC.grants[0] != 42 || len(coreRPC.apis) == 0 {
		t.Fatalf("catalog initialization not idempotent: menus=%d grants=%v apis=%d", len(coreRPC.menus), coreRPC.grants, len(coreRPC.apis))
	}
	coreRPC.fail = true
	if _, err := server.InitDatabase(ctx, &ops.Empty{}); err == nil {
		t.Fatal("Core RPC failure must not be reported as success")
	}
	if _, err := client.AccessProfile.Get(tenantCtx, profile.ID); err != nil {
		t.Fatalf("existing profile lost: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM other_module WHERE id = 7").Scan(&count); err != nil || count != 1 {
		t.Fatalf("shared table changed: count=%d err=%v", count, err)
	}
	var note sql.NullString
	if err := db.QueryRow("SELECT legacy_note FROM access_profiles").Scan(&note); err != nil {
		t.Fatalf("legacy column removed: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='legacy_profile_note'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("legacy index removed: count=%d err=%v", count, err)
	}
	otherCtx := hooks.SetTenantIDToContext(ctx, 2)
	if count, err := client.AccessProfile.Query().Count(otherCtx); err != nil || count != 0 {
		t.Fatalf("tenant isolation: count=%d err=%v", count, err)
	}
	if _, err := client.AccessProfile.Query().Count(ctx); err == nil {
		t.Fatal("query without tenant context must remain rejected")
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := server.InitDatabase(ctx, &ops.Empty{}); err == nil {
		t.Fatal("closed database must return initialization failure")
	}
}
