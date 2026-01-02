package session

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionBySessionIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionBySessionIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionBySessionIdLogic {
	return &GetSessionBySessionIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionBySessionIdLogic) GetSessionBySessionId(in *ops.SessionSIDReq) (*ops.SessionInfo, error) {
	// todo: add your logic here and delete this line

	return &ops.SessionInfo{}, nil
}
