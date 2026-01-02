package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProfileLogic {
	return &ListProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProfileLogic) ListProfile() (resp *types.AccessProfileListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
