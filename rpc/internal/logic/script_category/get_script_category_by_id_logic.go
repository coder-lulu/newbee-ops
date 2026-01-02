package script_category

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptCategoryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptCategoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptCategoryByIdLogic {
	return &GetScriptCategoryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptCategoryByIdLogic) GetScriptCategoryById(in *ops.IDReq) (*ops.ScriptCategoryInfo, error) {
	result, err := l.svcCtx.DB.ScriptCategory.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ScriptCategoryInfo{
		Id:          &result.ID,
		CreatedAt:   pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:   pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:      pointy.GetPointer(uint32(result.Status)),
		Name:        &result.Name,
		Code:        &result.Code,
		Description: &result.Description,
		ParentId:    &result.ParentID,
		SortOrder:   &result.SortOrder,
		Icon:        &result.Icon,
	}, nil
}
