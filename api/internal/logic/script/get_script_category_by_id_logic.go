package script

import (
	"context"
	"fmt"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptCategoryByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取脚本分类详情
func NewGetScriptCategoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptCategoryByIdLogic {
	return &GetScriptCategoryByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptCategoryByIdLogic) GetScriptCategoryById(req *types.IDFormReq) (resp *types.ScriptCategoryDetailResp, err error) {
	// Call RPC service
	rpcReq := &ops.IDReq{
		Id: req.Id,
	}

	result, err := l.svcCtx.OpsClient.GetScriptCategoryById(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to get script category by ID via RPC",
			logx.Field("id", req.Id),
			logx.Field("error", err))
		return &types.ScriptCategoryDetailResp{
			Code: 500,
			Msg:  fmt.Sprintf("获取脚本分类详情失败: %v", err),
			Data: types.ScriptCategoryInfo{},
		}, nil
	}

	// Convert RPC response to API response
	return &types.ScriptCategoryDetailResp{
		Code: 0,
		Msg:  "success",
		Data: types.ScriptCategoryInfo{
			Id:          getValue(result.Id),
			Name:        getValue(result.Name),
			Code:        getValue(result.Code),
			Description: getValue(result.Description),
			ParentId:    getValue(result.ParentId),
			SortOrder:   getValue(result.SortOrder),
			Icon:        getValue(result.Icon),
		},
	}, nil
}
