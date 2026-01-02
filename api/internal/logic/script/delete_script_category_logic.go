package script

import (
	"context"
	"fmt"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScriptCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除脚本分类
func NewDeleteScriptCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScriptCategoryLogic {
	return &DeleteScriptCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteScriptCategoryLogic) DeleteScriptCategory(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// Convert API types to RPC types
	rpcReq := &ops.IDsReq{
		Ids: req.Ids,
	}

	// Call RPC service
	result, err := l.svcCtx.OpsClient.DeleteScriptCategory(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to delete script category via RPC",
			logx.Field("ids", req.Ids),
			logx.Field("error", err))
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("删除脚本分类失败: %v", err),
		}, nil
	}

	// Return success response
	return &types.BaseResp{
		Code: uint32(result.Code),
		Msg:  result.Msg,
	}, nil
}
