package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProxyByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Proxy详情（需要JWT）
func NewGetProxyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyByIdLogic {
	return &GetProxyByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProxyByIdLogic) GetProxyById(req *types.IDReq) (resp *types.ProxyDetail, err error) {
	// todo: add your logic here and delete this line

	return
}
