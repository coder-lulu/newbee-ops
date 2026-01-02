package script_version

import (
	"context"

    "github.com/coder-lulu/newbee-ops-rpc/ent/scriptversion"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteScriptVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteScriptVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScriptVersionLogic {
	return &DeleteScriptVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteScriptVersionLogic) DeleteScriptVersion(in *ops.IDsReq) (*ops.BaseResp, error) {
	_, err := l.svcCtx.DB.ScriptVersion.Delete().Where(scriptversion.IDIn(in.Ids...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.DeleteSuccess }, nil
}
