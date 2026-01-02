package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PickProxyFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 从分组选择Proxy（需要JWT）
func NewPickProxyFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PickProxyFromGroupLogic {
	return &PickProxyFromGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PickProxyFromGroupLogic) PickProxyFromGroup(req *types.PickProxyFromGroupReq) (resp *types.ProxyPickResp, err error) {
	// todo: add your logic here and delete this line

	return
}
