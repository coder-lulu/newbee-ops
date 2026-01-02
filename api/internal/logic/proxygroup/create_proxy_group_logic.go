package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProxyGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建Proxy分组（需要JWT）
func NewCreateProxyGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProxyGroupLogic {
	return &CreateProxyGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProxyGroupLogic) CreateProxyGroup(req *types.ProxyGroupInfo) (resp *types.BaseIDResp, err error) {
	// todo: add your logic here and delete this line

	return
}
