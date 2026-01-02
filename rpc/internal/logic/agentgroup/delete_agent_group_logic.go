package agentgroup

import (
	"context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/agentgroup"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentGroupLogic {
	return &DeleteAgentGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAgentGroupLogic) DeleteAgentGroup(in *ops.IDsReq) (*ops.BaseResp, error) {
	_, err := l.svcCtx.DB.AgentGroup.Delete().Where(agentgroup.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.DeleteSuccess }, nil
}
