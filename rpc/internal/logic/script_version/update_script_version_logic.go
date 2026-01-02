package script_version

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScriptVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScriptVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptVersionLogic {
	return &UpdateScriptVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScriptVersionLogic) UpdateScriptVersion(in *ops.ScriptVersionInfo) (*ops.BaseResp, error) {
	err:= l.svcCtx.DB.ScriptVersion.UpdateOneID(*in.Id).
			SetNotNilScriptID(in.ScriptId).
			SetNotNilVersion(in.Version).
			SetNotNilContent(in.Content).
			SetParameters(unmarshalJSONMap(in.Parameters)).
			SetNotNilScriptType(in.ScriptType).
			SetNotNilExecutor(in.Executor).
			SetNotNilChangeLog(in.ChangeLog).
			SetNotNilCreatedBy(in.CreatedBy).
			SetNotNilChecksum(in.Checksum).
			Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
