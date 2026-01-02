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

type CreateScriptCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateScriptCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptCategoryLogic {
	return &CreateScriptCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateScriptCategoryLogic) CreateScriptCategory(in *ops.ScriptCategoryInfo) (*ops.BaseIDResp, error) {
	if in.Name == nil || *in.Name == "" {
		return nil, newInvalidArgumentError("category name is required")
	}
	if in.Code == nil || *in.Code == "" {
		return nil, newInvalidArgumentError("category code is required")
	}

	query := l.svcCtx.DB.ScriptCategory.Create().
		SetNotNilName(in.Name).
		SetNotNilCode(in.Code).
		SetNotNilDescription(in.Description).
		SetNotNilParentID(in.ParentId).
		SetNotNilSortOrder(in.SortOrder).
		SetNotNilIcon(in.Icon)

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}

	result, err := query.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess}, nil
}
