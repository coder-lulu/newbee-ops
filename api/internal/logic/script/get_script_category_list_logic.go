package script

import (
	"context"
	"fmt"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取脚本分类列表
func NewGetScriptCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptCategoryListLogic {
	return &GetScriptCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptCategoryListLogic) GetScriptCategoryList(req *types.ScriptCategoryListReq) (resp *types.ScriptCategoryListResp, err error) {
	// Build RPC request - get all categories (Page=1, PageSize=10000 for tree mode)
	rpcReq := &ops.ScriptCategoryListReq{
		Page:     1,
		PageSize: 10000, // Get all for tree building
	}

	// If parentId is specified, filter by it
	if req.ParentId != 0 {
		rpcReq.ParentId = pointy(req.ParentId)
	}

	// Call RPC service
	result, err := l.svcCtx.OpsClient.GetScriptCategoryList(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to get script category list via RPC",
			logx.Field("parentId", req.ParentId),
			logx.Field("tree", req.Tree),
			logx.Field("error", err))
		return nil, fmt.Errorf("获取脚本分类列表失败: %v", err)
	}

	// Convert RPC response to API response
	var categoryItems []types.ScriptCategoryItem

	if req.Tree {
		// Build tree structure
		categoryItems = l.buildCategoryTree(result.Data)
	} else {
		// Return flat list
		categoryItems = make([]types.ScriptCategoryItem, 0, len(result.Data))
		for _, item := range result.Data {
			categoryItems = append(categoryItems, convertRPCInfoToAPI(item))
		}
	}

	return &types.ScriptCategoryListResp{
		Code: 0,
		Msg:  "success",
		Data: categoryItems,
	}, nil
}

// buildCategoryTree builds a tree structure from flat list
func (l *GetScriptCategoryListLogic) buildCategoryTree(flatList []*ops.ScriptCategoryInfo) []types.ScriptCategoryItem {
	// Map to store all categories by ID
	categoryMap := make(map[uint64]*types.ScriptCategoryItem)

	// First pass: convert all items
	for _, item := range flatList {
		apiItem := convertRPCInfoToAPI(item)
		categoryMap[apiItem.Id] = &apiItem
	}

	// Second pass: build parent-child relationships
	for _, item := range categoryMap {
		if item.ParentId != 0 {
			// Child category - add to parent
			if parent, exists := categoryMap[item.ParentId]; exists {
				if parent.Children == nil {
					parent.Children = make([]types.ScriptCategoryItem, 0)
				}
				parent.Children = append(parent.Children, *item)
			}
		}
	}

	// Third pass: collect root categories (after tree is fully built)
	var rootCategories []types.ScriptCategoryItem
	for _, item := range categoryMap {
		if item.ParentId == 0 {
			rootCategories = append(rootCategories, *item)
		}
	}

	return rootCategories
}

// convertRPCInfoToAPI converts RPC ScriptCategoryInfo to API ScriptCategoryItem
func convertRPCInfoToAPI(rpcInfo *ops.ScriptCategoryInfo) types.ScriptCategoryItem {
	return types.ScriptCategoryItem{
		Id:          getValue(rpcInfo.Id),
		Name:        getValue(rpcInfo.Name),
		Code:        getValue(rpcInfo.Code),
		Description: getValue(rpcInfo.Description),
		ParentId:    getValue(rpcInfo.ParentId),
		SortOrder:   getValue(rpcInfo.SortOrder),
		Icon:        getValue(rpcInfo.Icon),
		ScriptCount: 0, // RPC doesn't provide this, would need separate query
	}
}
