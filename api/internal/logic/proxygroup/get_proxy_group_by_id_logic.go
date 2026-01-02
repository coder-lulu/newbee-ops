package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy分组详情（需要JWT）
func NewGetProxyGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupByIdLogic {
	return &GetProxyGroupByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyGroupByIdLogic) GetProxyGroupById(req *types.IDPathReq) (resp *types.ProxyGroupDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
