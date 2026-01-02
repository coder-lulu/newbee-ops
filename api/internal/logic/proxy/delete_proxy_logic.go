package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Proxy（需要JWT）
func NewDeleteProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProxyLogic {
	return &DeleteProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProxyLogic) DeleteProxy(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
