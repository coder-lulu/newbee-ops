package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy分组列表（需要JWT）
func NewGetProxyGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupListLogic {
	return &GetProxyGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyGroupListLogic) GetProxyGroupList(req *types.ProxyGroupListReq) (resp *types.ProxyGroupListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
