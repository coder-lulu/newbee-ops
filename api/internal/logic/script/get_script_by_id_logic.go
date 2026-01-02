package script

import (
	"context"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取脚本详情
func NewGetScriptByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptByIdLogic {
	return &GetScriptByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScriptByIdLogic) GetScriptById(req *types.IDFormReq) (resp *types.ScriptDetailResp, err error) {
	// 调用RPC服务
	rpcReq := &ops.IDReq{
		Id: req.Id,
	}

	result, err := l.svcCtx.OpsClient.GetScriptById(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to get script by id via RPC",
			logx.Field("error", err),
			logx.Field("id", req.Id))
		return nil, err
	}

	// 转换响应 - 反序列化JSON字段
	scriptDetail := &types.ScriptDetail{
		Id:                   getValue(result.Id),
		CreatedAt:            getValue(result.CreatedAt),
		UpdatedAt:            getValue(result.UpdatedAt),
		Status:               getValue(result.Status),
		Name:                 getValue(result.Name),
		Code:                 getValue(result.Code),
		Description:          getValue(result.Description),
		CategoryId:           getValue(result.CategoryId),
		Tags:                 unmarshalJSONStringSlice(result.Tags),
		ScriptType:           getValue(result.ScriptType),
		Executor:             getValue(result.Executor),
		Content:              getValue(result.Content),
		IsTemplate:           getValue(result.IsTemplate),
		TemplateEngine:       getValue(result.TemplateEngine),
		Parameters:           unmarshalJSONMap(result.Parameters),
		Version:              getValue(result.Version),
		BaseVersionId:        getValue(result.BaseVersionId),
		IsLatest:             getValue(result.IsLatest),
		TargetCiTypes:        unmarshalJSONStringSlice(result.TargetCiTypes),
		TargetOsTypes:        unmarshalJSONStringSlice(result.TargetOsTypes),
		TargetSelector:       unmarshalJSONMap(result.TargetSelector),
		CredentialRef:        getValue(result.CredentialRef),
		RequiredCapabilities: unmarshalJSONStringSlice(result.RequiredCapabilities),
		DefaultTimeout:       getValue(result.DefaultTimeout),
		DefaultWorkdir:       getValue(result.DefaultWorkdir),
		DefaultEnv:           unmarshalJSONMap(result.DefaultEnv),
		RequireConfirmation:  getValue(result.RequireConfirmation),
		RiskLevel:            getValue(result.RiskLevel),
		Schedulable:          getValue(result.Schedulable),
		DefaultSchedule:      getValue(result.DefaultSchedule),
		ExecutionCount:       getValue(result.ExecutionCount),
		SuccessCount:         getValue(result.SuccessCount),
		FailureCount:         getValue(result.FailureCount),
		LastExecutedAt:       getValue(result.LastExecutedAt),
	}

	return &types.ScriptDetailResp{
		Code: 0,
		Msg:  "success",
		Data: *scriptDetail,
	}, nil
}
