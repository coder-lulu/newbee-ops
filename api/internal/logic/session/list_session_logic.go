package session

import (
    "context"
    "github.com/coder-lulu/newbee-ops-api/internal/svc"
)

type ListSessionLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewListSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionLogic {
    return &ListSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListSessionLogic) List(status, ciId string, page, size int) ([]*svc.Session, int) {
    tenant := ""
    if l.svcCtx.ContextManager != nil { tenant = l.svcCtx.ContextManager.GetTenantID(l.ctx) }
    return l.svcCtx.SessionStore.ListBy(tenant, status, ciId, page, size)
}

