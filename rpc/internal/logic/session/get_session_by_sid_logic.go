package session

import (
    "context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/session"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetSessionBySIDLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewGetSessionBySIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionBySIDLogic {
    return &GetSessionBySIDLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetSessionBySIDLogic) GetSessionBySessionId(in *newbee_ops_rpc.SessionSIDReq) (*newbee_ops_rpc.SessionInfo, error) {
    v, err := l.svcCtx.DB.Session.Query().Where(session.SessionIDEQ(in.GetSessionId())).Only(l.ctx)
    if err != nil { return nil, err }
    ca := v.CreatedAt.Unix()
    ua := v.UpdatedAt.Unix()
    return &newbee_ops_rpc.SessionInfo{
        Id: &v.ID,
        CreatedAt: &ca,
        UpdatedAt: &ua,
        SessionId: &v.SessionID,
        UserId: &v.UserID,
        CiId: &v.CiID,
        Protocol: &v.Protocol,
        ProxyId: &v.ProxyID,
        Endpoint: &v.Endpoint,
        ExpiresAt: &v.ExpiresAt,
        StatusStr: &v.StatusStr,
        ClosedAt: &v.ClosedAt,
    }, nil
}
