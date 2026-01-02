package session

import (
    "context"
    "time"

    "github.com/coder-lulu/newbee-ops-rpc/types/ops"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
)

type CloseSessionLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCloseSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseSessionLogic {
    return &CloseSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CloseSessionLogic) Close(sessionId string) (ok bool, closedAt int64, reason string) {
    tenant := ""
    if l.svcCtx.ContextManager != nil { tenant = l.svcCtx.ContextManager.GetTenantID(l.ctx) }
    ts := time.Now().Unix()
    ok, reason = l.svcCtx.SessionStore.Close(sessionId, tenant, ts)
    if ok { return true, ts, "" }
    // 如果内存未命中，尝试 RPC 更新
    if l.svcCtx.OpsClient != nil {
        // 查询会话ID
        if info, err := l.svcCtx.OpsClient.GetSessionBySessionId(l.ctx, &ops.SessionSIDReq{SessionId: sessionId}); err == nil && info != nil && info.Id != nil {
            _, _ = l.svcCtx.OpsClient.UpdateSession(l.ctx, &ops.SessionInfo{Id: info.Id, StatusStr: strPtr("closed"), ClosedAt: &ts})
            return true, ts, ""
        }
    }
    return false, 0, reason
}

func strPtr(s string) *string { return &s }
