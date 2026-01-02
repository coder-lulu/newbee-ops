package script_category

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/ent/scriptcategory"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptCategoryListLogic {
	return &GetScriptCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptCategoryListLogic) GetScriptCategoryList(in *ops.ScriptCategoryListReq) (*ops.ScriptCategoryListResp, error) {
	var predicates []predicate.ScriptCategory

	if in.Name != nil {
		predicates = append(predicates, scriptcategory.NameContains(*in.Name))
	}
	if in.Code != nil {
		predicates = append(predicates, scriptcategory.CodeEQ(*in.Code))
	}
	if in.ParentId != nil {
		predicates = append(predicates, scriptcategory.ParentIDEQ(*in.ParentId))
	}
	if in.Status != nil {
		predicates = append(predicates, scriptcategory.StatusEQ(uint8(*in.Status)))
	}

	pageResult, err := l.svcCtx.DB.ScriptCategory.Query().
		Where(predicates...).
		Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ScriptCategoryListResp{}
	resp.Total = pageResult.PageDetails.Total

	for _, v := range pageResult.List {
		resp.Data = append(resp.Data, &ops.ScriptCategoryInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:      pointy.GetPointer(uint32(v.Status)),
			Name:        &v.Name,
			Code:        &v.Code,
			Description: &v.Description,
			ParentId:    &v.ParentID,
			SortOrder:   &v.SortOrder,
			Icon:        &v.Icon,
		})
	}

	return resp, nil
}
