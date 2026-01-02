package script

import (
	"context"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取脚本列表
func NewGetScriptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptListLogic {
	return &GetScriptListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptListLogic) GetScriptList(req *types.ScriptListReq) (resp *types.ScriptListResp, err error) {
	// 构建RPC请求
	rpcReq := &ops.ScriptListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 可选参数
	if req.Name != "" {
		rpcReq.Name = &req.Name
	}
	if req.Code != "" {
		rpcReq.Code = &req.Code
	}
	if req.CategoryId > 0 {
		rpcReq.CategoryId = &req.CategoryId
	}
	if req.ScriptType != "" {
		rpcReq.ScriptType = &req.ScriptType
	}
	if req.Executor != "" {
		rpcReq.Executor = &req.Executor
	}
	if req.RiskLevel != "" {
		rpcReq.RiskLevel = &req.RiskLevel
	}

	// 调用RPC服务
	result, err := l.svcCtx.OpsClient.GetScriptList(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to get script list via RPC",
			logx.Field("error", err))
		return nil, err
	}

	// 转换响应
	resp = &types.ScriptListResp{
		Code: 0,
		Msg:  "success",
		Data: types.ScriptListData{
			Total: result.Total,
			Data:  make([]types.ScriptItem, 0),
		},
	}

	for _, item := range result.Data {
		resp.Data.Data = append(resp.Data.Data, types.ScriptItem{
			Id:             getValue(item.Id),
			CreatedAt:      getValue(item.CreatedAt),
			Name:           getValue(item.Name),
			Code:           getValue(item.Code),
			Description:    getValue(item.Description),
			CategoryId:     getValue(item.CategoryId),
			ScriptType:     getValue(item.ScriptType),
			Executor:       getValue(item.Executor),
			IsTemplate:     getValue(item.IsTemplate),
			Version:        getValue(item.Version),
			IsLatest:       getValue(item.IsLatest),
			RiskLevel:      getValue(item.RiskLevel),
			Schedulable:    getValue(item.Schedulable),
			ExecutionCount: getValue(item.ExecutionCount),
			SuccessCount:   getValue(item.SuccessCount),
			FailureCount:   getValue(item.FailureCount),
			LastExecutedAt: getValue(item.LastExecutedAt),
		})
	}

	return resp, nil
}
