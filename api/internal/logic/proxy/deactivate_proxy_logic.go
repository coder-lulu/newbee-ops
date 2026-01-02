package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeactivateProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手动下线Proxy（需要JWT）
func NewDeactivateProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateProxyLogic {
	return &DeactivateProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivateProxyLogic) DeactivateProxy(req *types.IDReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
