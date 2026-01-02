package script

import (
	"context"
	"fmt"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScriptCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建脚本分类
func NewCreateScriptCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptCategoryLogic {
	return &CreateScriptCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateScriptCategoryLogic) CreateScriptCategory(req *types.ScriptCategoryInfo) (resp *types.BaseIDResp, err error) {
	// Convert API types to RPC types (add pointers)
	rpcReq := &ops.ScriptCategoryInfo{
		Name:        &req.Name,
		Code:        &req.Code,
		Description: &req.Description,
		ParentId:    pointy(req.ParentId),
		SortOrder:   pointy(req.SortOrder),
		Icon:        &req.Icon,
	}

	// Call RPC service
	result, err := l.svcCtx.OpsClient.CreateScriptCategory(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to create script category via RPC",
			logx.Field("name", req.Name),
			logx.Field("code", req.Code),
			logx.Field("error", err))
		return &types.BaseIDResp{
			Code: 500,
			Msg:  fmt.Sprintf("创建脚本分类失败: %v", err),
			Id:   0,
		}, nil
	}

	// Return success response
	return &types.BaseIDResp{
		Code: 0, // RPC BaseIDResp doesn't have Code, 0 means success
		Msg:  result.Msg,
		Id:   result.Id,
	}, nil
}
