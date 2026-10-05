package session

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/coder-lulu/newbee-ops-rpc/ent"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

func makeSvcCtxForTest(t *testing.T) *svc.ServiceContext {
	t.Helper()
	connection, err := sql.Open("sqlite", "file:session-crud?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	connection.SetMaxOpenConns(1)
	db := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, connection)))
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Schema.Create(context.Background()); err != nil {
		t.Fatal(err)
	}
	return &svc.ServiceContext{DB: db}
}

func TestSessionCRUD_Filter(t *testing.T) {
	ctx := context.Background()
	svcCtx := makeSvcCtxForTest(t)

	// Create
	sid := "sess-001"
	user := "u1"
	ci := "ci-1"
	proto := "ssh"
	status := "active"
	_, err := NewCreateSessionLogic(ctx, svcCtx).CreateSession(&pb.SessionInfo{SessionId: &sid, UserId: &user, CiId: &ci, Protocol: &proto, StatusStr: &status})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}

	// Get by SID
	g1, err := NewGetSessionBySIDLogic(ctx, svcCtx).GetSessionBySessionId(&pb.SessionSIDReq{SessionId: sid})
	if err != nil {
		t.Fatalf("get by sid err: %v", err)
	}
	if g1.GetSessionId() != sid {
		t.Fatalf("sid mismatch")
	}

	// Update to closed
	closed := "closed"
	_, err = NewUpdateSessionLogic(ctx, svcCtx).UpdateSession(&pb.SessionInfo{Id: g1.Id, StatusStr: &closed})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}

	// List filter
	lresp, err := NewGetSessionListLogic(ctx, svcCtx).GetSessionList(&pb.SessionListReq{Page: 1, PageSize: 10, StatusStr: &closed})
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if lresp.GetTotal() == 0 {
		t.Fatalf("list empty after update")
	}

	// Delete and verify the filtered list no longer contains the session.
	if _, err := NewDeleteSessionLogic(ctx, svcCtx).DeleteSession(&pb.IDsReq{Ids: []uint64{g1.GetId()}}); err != nil {
		t.Fatalf("delete err: %v", err)
	}
	lresp, err = NewGetSessionListLogic(ctx, svcCtx).GetSessionList(&pb.SessionListReq{Page: 1, PageSize: 10, StatusStr: &closed})
	if err != nil {
		t.Fatalf("list after delete err: %v", err)
	}
	if lresp.GetTotal() != 0 || len(lresp.GetData()) != 0 {
		t.Fatal("deleted session remains in list")
	}
}
