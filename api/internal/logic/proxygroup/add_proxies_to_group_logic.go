package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProxiesToGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加Proxy到分组（需要JWT）
func NewAddProxiesToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProxiesToGroupLogic {
	return &AddProxiesToGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProxiesToGroupLogic) AddProxiesToGroup(req *types.AddProxiesToGroupReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
