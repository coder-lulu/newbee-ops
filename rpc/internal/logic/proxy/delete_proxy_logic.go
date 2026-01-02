package proxy

import (
	"context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteProxyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProxyLogic {
	return &DeleteProxyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProxyLogic) DeleteProxy(in *ops.IDsReq) (*ops.BaseResp, error) {
	_, err := l.svcCtx.DB.Proxy.Delete().Where(proxy.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.DeleteSuccess }, nil
}
