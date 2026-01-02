package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy列表（需要JWT）
func NewGetProxyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyListLogic {
	return &GetProxyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyListLogic) GetProxyList(req *types.ProxyListReq) (resp *types.ProxyListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
