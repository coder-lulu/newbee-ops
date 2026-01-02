package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProxiesFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 从分组移除Proxy（需要JWT）
func NewRemoveProxiesFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProxiesFromGroupLogic {
	return &RemoveProxiesFromGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveProxiesFromGroupLogic) RemoveProxiesFromGroup(req *types.RemoveProxiesFromGroupReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
