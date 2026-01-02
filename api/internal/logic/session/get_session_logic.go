package session

import (
    "context"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
)

type GetSessionLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
    return &GetSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetSessionLogic) Get(id string) (*svc.Session, bool) {
    sess, ok := l.svcCtx.SessionStore.Get(id)
    if !ok { return nil, false }
    tenant := ""
    if l.svcCtx.ContextManager != nil { tenant = l.svcCtx.ContextManager.GetTenantID(l.ctx) }
    if tenant != "" && sess.TenantId != "" && tenant != sess.TenantId { return nil, false }
    return sess, true
}

