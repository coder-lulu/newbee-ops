package script_category

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/msg/errormsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScriptCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScriptCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptCategoryLogic {
	return &UpdateScriptCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScriptCategoryLogic) UpdateScriptCategory(in *ops.ScriptCategoryInfo) (*ops.BaseResp, error) {
	if in.Id == nil {
		return nil, newInvalidArgumentError("category id is required")
	}

	query := l.svcCtx.DB.ScriptCategory.UpdateOneID(*in.Id).
		SetNotNilName(in.Name).
		SetNotNilCode(in.Code).
		SetNotNilDescription(in.Description).
		SetNotNilParentID(in.ParentId).
		SetNotNilSortOrder(in.SortOrder).
		SetNotNilIcon(in.Icon)

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}

	err := query.Exec(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.BaseResp{Msg: errormsg.UpdateSuccess}, nil
}
