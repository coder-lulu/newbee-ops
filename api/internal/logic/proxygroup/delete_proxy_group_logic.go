package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProxyGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Proxy分组（需要JWT）
func NewDeleteProxyGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProxyGroupLogic {
	return &DeleteProxyGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProxyGroupLogic) DeleteProxyGroup(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
