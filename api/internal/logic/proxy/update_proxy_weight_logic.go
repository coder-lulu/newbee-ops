package proxy

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyWeightLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新Proxy权重（需要JWT）
func NewUpdateProxyWeightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyWeightLogic {
	return &UpdateProxyWeightLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProxyWeightLogic) UpdateProxyWeight(req *types.UpdateProxyWeightReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
