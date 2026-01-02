package script_version

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScriptVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateScriptVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptVersionLogic {
	return &CreateScriptVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateScriptVersionLogic) CreateScriptVersion(in *ops.ScriptVersionInfo) (*ops.BaseIDResp, error) {
    result, err := l.svcCtx.DB.ScriptVersion.Create().
			SetNotNilScriptID(in.ScriptId).
			SetNotNilVersion(in.Version).
			SetNotNilContent(in.Content).
			SetParameters(unmarshalJSONMap(in.Parameters)).
			SetNotNilScriptType(in.ScriptType).
			// SetNotNilExecutor removed - field does not exist.
			SetNotNilChangeLog(in.ChangeLog).
			SetNotNilCreatedBy(in.CreatedBy).
			SetNotNilChecksum(in.Checksum).
			Save(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess }, nil
}
