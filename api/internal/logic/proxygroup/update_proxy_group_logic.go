package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新Proxy分组（需要JWT）
func NewUpdateProxyGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyGroupLogic {
	return &UpdateProxyGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProxyGroupLogic) UpdateProxyGroup(req *types.ProxyGroupInfo) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
