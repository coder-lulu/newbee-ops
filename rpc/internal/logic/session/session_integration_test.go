package session

import (
    "context"
    "testing"

    _ "modernc.org/sqlite"

    "github.com/coder-lulu/newbee-ops-rpc/ent/enttest"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

func makeSvcCtxForTest(t *testing.T) *svc.ServiceContext {
    t.Helper()
    db := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
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
    if err != nil { t.Fatalf("create err: %v", err) }

    // Get by SID
    g1, err := NewGetSessionBySIDLogic(ctx, svcCtx).GetSessionBySessionId(&pb.SessionSIDReq{SessionId: sid})
    if err != nil { t.Fatalf("get by sid err: %v", err) }
    if g1.GetSessionId() != sid { t.Fatalf("sid mismatch") }

    // Update to closed
    closed := "closed"
    _, err = NewUpdateSessionLogic(ctx, svcCtx).UpdateSession(&pb.SessionInfo{Id: g1.Id, StatusStr: &closed})
    if err != nil { t.Fatalf("update err: %v", err) }

    // List filter
    lresp, err := NewGetSessionListLogic(ctx, svcCtx).GetSessionList(&pb.SessionListReq{Page: 1, PageSize: 10, StatusStr: &closed})
    if err != nil { t.Fatalf("list err: %v", err) }
    if lresp.GetTotal() == 0 { t.Fatalf("list empty after update") }
}

