package script

import (
	"context"
	"fmt"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScriptCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新脚本分类
func NewUpdateScriptCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptCategoryLogic {
	return &UpdateScriptCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateScriptCategoryLogic) UpdateScriptCategory(req *types.ScriptCategoryInfo) (resp *types.BaseResp, err error) {
	// Convert API types to RPC types (add pointers)
	rpcReq := &ops.ScriptCategoryInfo{
		Id:          &req.Id,
		Name:        &req.Name,
		Code:        &req.Code,
		Description: &req.Description,
		ParentId:    pointy(req.ParentId),
		SortOrder:   pointy(req.SortOrder),
		Icon:        &req.Icon,
	}

	// Call RPC service
	result, err := l.svcCtx.OpsClient.UpdateScriptCategory(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to update script category via RPC",
			logx.Field("id", req.Id),
			logx.Field("name", req.Name),
			logx.Field("error", err))
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("更新脚本分类失败: %v", err),
		}, nil
	}

	// Return success response
	return &types.BaseResp{
		Code: uint32(result.Code),
		Msg:  result.Msg,
	}, nil
}
