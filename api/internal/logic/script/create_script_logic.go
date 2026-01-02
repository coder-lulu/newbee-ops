package script

import (
	"context"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScriptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建脚本
func NewCreateScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptLogic {
	return &CreateScriptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateScriptLogic) CreateScript(req *types.ScriptInfo) (resp *types.BaseIDResp, err error) {
	// 构建RPC请求
	rpcReq := &ops.ScriptInfo{
		Name:                 &req.Name,
		Code:                 &req.Code,
		ScriptType:           &req.ScriptType,
		Executor:             &req.Executor,
		Content:              &req.Content,
	}

	// 可选的基础字段
	if req.Description != "" {
		rpcReq.Description = &req.Description
	}
	if req.CategoryId > 0 {
		rpcReq.CategoryId = &req.CategoryId
	}

	// JSON数组字段需要序列化
	if len(req.Tags) > 0 {
		rpcReq.Tags = marshalJSONStringSlice(req.Tags)
	}
	if len(req.TargetCiTypes) > 0 {
		rpcReq.TargetCiTypes = marshalJSONStringSlice(req.TargetCiTypes)
	}
	if len(req.TargetOsTypes) > 0 {
		rpcReq.TargetOsTypes = marshalJSONStringSlice(req.TargetOsTypes)
	}
	if len(req.RequiredCapabilities) > 0 {
		rpcReq.RequiredCapabilities = marshalJSONStringSlice(req.RequiredCapabilities)
	}

	// JSON Map字段需要序列化
	if len(req.Parameters) > 0 {
		rpcReq.Parameters = marshalJSONMap(req.Parameters)
	}
	if len(req.TargetSelector) > 0 {
		rpcReq.TargetSelector = marshalJSONMap(req.TargetSelector)
	}
	if len(req.DefaultEnv) > 0 {
		rpcReq.DefaultEnv = marshalJSONMap(req.DefaultEnv)
	}

	// 模板相关字段
	rpcReq.IsTemplate = &req.IsTemplate
	if req.TemplateEngine != "" {
		rpcReq.TemplateEngine = &req.TemplateEngine
	}

	// 凭证和执行配置
	if req.CredentialRef != "" {
		rpcReq.CredentialRef = &req.CredentialRef
	}
	if req.DefaultTimeout > 0 {
		rpcReq.DefaultTimeout = &req.DefaultTimeout
	}
	if req.DefaultWorkdir != "" {
		rpcReq.DefaultWorkdir = &req.DefaultWorkdir
	}

	// 安全控制
	rpcReq.RequireConfirmation = &req.RequireConfirmation
	if req.RiskLevel != "" {
		rpcReq.RiskLevel = &req.RiskLevel
	}

	// 调度配置
	rpcReq.Schedulable = &req.Schedulable
	if req.DefaultSchedule != "" {
		rpcReq.DefaultSchedule = &req.DefaultSchedule
	}

	// 调用RPC服务
	result, err := l.svcCtx.OpsClient.CreateScript(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorw("Failed to create script via RPC",
			logx.Field("error", err))
		return nil, err
	}

	// 返回结果
	return &types.BaseIDResp{
		Code: 0,
		Msg:  result.Msg,
		Id:   result.Id,
	}, nil
}
