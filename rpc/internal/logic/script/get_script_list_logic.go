package script

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/ent/script"
	"github.com/coder-lulu/newbee-ops-rpc/ent/scriptcategory"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptListLogic {
	return &GetScriptListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptListLogic) GetScriptList(in *ops.ScriptListReq) (*ops.ScriptListResp, error) {
	var predicates []predicate.Script

	// 支持按名称模糊查询
	if in.Name != nil {
		predicates = append(predicates, script.NameContains(*in.Name))
	}

	// 支持按代码精确查询
	if in.Code != nil {
		predicates = append(predicates, script.CodeEQ(*in.Code))
	}

	// 支持按分类ID查询（包含子分类）
	if in.CategoryId != nil {
		// 递归查询该分类及其所有子分类的ID
		categoryIds, err := l.getAllCategoryIds(*in.CategoryId)
		if err != nil {
			l.Logger.Errorw("Failed to get category ids",
				logx.Field("category_id", *in.CategoryId),
				logx.Field("error", err))
			// 如果查询失败，回退到只查询当前分类
			predicates = append(predicates, script.CategoryIDEQ(*in.CategoryId))
		} else {
			// 使用IN查询所有分类ID（包含父分类和所有子分类）
			predicates = append(predicates, script.CategoryIDIn(categoryIds...))
		}
	}

	// 支持按脚本类型查询
	if in.ScriptType != nil {
		predicates = append(predicates, script.ScriptTypeEQ(*in.ScriptType))
	}

	// 支持按执行器类型查询
	if in.Executor != nil {
		predicates = append(predicates, script.ExecutorEQ(*in.Executor))
	}

	// 支持按风险级别查询
	if in.RiskLevel != nil {
		predicates = append(predicates, script.RiskLevelEQ(*in.RiskLevel))
	}

	// 支持按状态查询
	if in.Status != nil {
		predicates = append(predicates, script.StatusEQ(uint8(*in.Status)))
	}

	pageResult, err := l.svcCtx.DB.Script.Query().
		Where(predicates...).
		Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ScriptListResp{}
	resp.Total = pageResult.PageDetails.Total

	for _, v := range pageResult.List {
		resp.Data = append(resp.Data, &ops.ScriptInfo{
			Id:                   &v.ID,
			CreatedAt:            pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:            pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:               pointy.GetPointer(uint32(v.Status)),
			Name:                 &v.Name,
			Code:                 &v.Code,
			Description:          &v.Description,
			CategoryId:           &v.CategoryID,
			Tags: pointy.GetPointer(mustMarshalJSON(v.Tags)),
			ScriptType:           &v.ScriptType,
			Executor:             &v.Executor,
			Content:              &v.Content,
			IsTemplate:           &v.IsTemplate,
			TemplateEngine:       &v.TemplateEngine,
			Parameters: pointy.GetPointer(mustMarshalJSON(v.Parameters)),
			Version:              &v.Version,
			BaseVersionId:        &v.BaseVersionID,
			IsLatest:             &v.IsLatest,
			TargetCiTypes: pointy.GetPointer(mustMarshalJSON(v.TargetCiTypes)),
			TargetOsTypes: pointy.GetPointer(mustMarshalJSON(v.TargetOsTypes)),
			TargetSelector: pointy.GetPointer(mustMarshalJSON(v.TargetSelector)),
			CredentialRef:        &v.CredentialRef,
			RequiredCapabilities: pointy.GetPointer(mustMarshalJSON(v.RequiredCapabilities)),
			DefaultTimeout:       &v.DefaultTimeout,
			DefaultWorkdir:       &v.DefaultWorkdir,
			DefaultEnv: pointy.GetPointer(mustMarshalJSON(v.DefaultEnv)),
			RequireConfirmation:  &v.RequireConfirmation,
			RiskLevel:            &v.RiskLevel,
			Schedulable:          &v.Schedulable,
			DefaultSchedule:      &v.DefaultSchedule,
			ExecutionCount:       &v.ExecutionCount,
			SuccessCount:         &v.SuccessCount,
			FailureCount:         &v.FailureCount,
			LastExecutedAt:       &v.LastExecutedAt,
		})
	}

	return resp, nil
}

// getAllCategoryIds 递归获取指定分类及其所有子分类的ID列表
func (l *GetScriptListLogic) getAllCategoryIds(categoryId uint64) ([]uint64, error) {
	// 结果集合，包含父分类ID本身
	categoryIds := []uint64{categoryId}

	// 递归查询所有子分类
	children, err := l.getChildCategoryIds(categoryId)
	if err != nil {
		return nil, err
	}

	categoryIds = append(categoryIds, children...)
	return categoryIds, nil
}

// getChildCategoryIds 递归查询所有子分类ID
func (l *GetScriptListLogic) getChildCategoryIds(parentId uint64) ([]uint64, error) {
	var result []uint64

	// 查询直接子分类
	directChildren, err := l.svcCtx.DB.ScriptCategory.Query().
		Where(scriptcategory.ParentIDEQ(parentId)).
		Select(scriptcategory.FieldID).
		All(l.ctx)

	if err != nil {
		return nil, err
	}

	// 如果没有子分类，返回空列表
	if len(directChildren) == 0 {
		return result, nil
	}

	// 遍历直接子分类，递归查询其子分类
	for _, child := range directChildren {
		// 添加当前子分类ID
		result = append(result, child.ID)

		// 递归查询该子分类的子分类
		grandChildren, err := l.getChildCategoryIds(child.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, grandChildren...)
	}

	return result, nil
}
