package session

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/coder-lulu/newbee-ops-api/internal/config"
    svcpkg "github.com/coder-lulu/newbee-ops-api/internal/svc"
)

func TestCloseSessionHandler(t *testing.T) {
    _ = config.Config{}
    ctx := &svcpkg.ServiceContext{SessionStore: svcpkg.NewSessionStore()}
    // prepopulate a session
    ctx.SessionStore.Start(&svcpkg.Session{ID: "sid-1", TenantId: "", Status: "active"})

    body := map[string]any{"sessionId": "sid-1"}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/ops/session/close", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    Close(ctx)(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("unexpected status: %d, body=%s", w.Code, w.Body.String())
    }
    var resp struct { Ok bool `json:"ok"` }
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("bad json: %v", err)
    }
    if !resp.Ok { t.Fatalf("expected ok true, got body=%s", w.Body.String()) }
}
