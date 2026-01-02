package proxy_group_member

import (
	"context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/proxygroupmember"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteProxyGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProxyGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProxyGroupMemberLogic {
	return &DeleteProxyGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProxyGroupMemberLogic) DeleteProxyGroupMember(in *ops.IDsReq) (*ops.BaseResp, error) {
	_, err := l.svcCtx.DB.ProxyGroupMember.Delete().Where(proxygroupmember.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.DeleteSuccess }, nil
}
