package script_version

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptVersionByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptVersionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptVersionByIdLogic {
	return &GetScriptVersionByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptVersionByIdLogic) GetScriptVersionById(in *ops.IDReq) (*ops.ScriptVersionInfo, error) {
	result, err := l.svcCtx.DB.ScriptVersion.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ScriptVersionInfo{
		Id:          &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		ScriptId:	&result.ScriptID,
		Version:	&result.Version,
		Content:	&result.Content,
		Parameters: pointy.GetPointer(mustMarshalJSON(result.Parameters)),
		ScriptType:	&result.ScriptType,
		Executor:	&result.Executor,
		ChangeLog:	&result.ChangeLog,
		CreatedBy:	&result.CreatedBy,
		Checksum:	&result.Checksum,
	}, nil
}

