package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProfileLogic {
	return &DeleteProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProfileLogic) DeleteProfile() error {
	// todo: add your logic here and delete this line

	return nil
}
