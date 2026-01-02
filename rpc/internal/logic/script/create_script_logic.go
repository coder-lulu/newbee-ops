package script

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/msg/errormsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScriptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptLogic {
	return &CreateScriptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateScriptLogic) CreateScript(in *ops.ScriptInfo) (*ops.BaseIDResp, error) {
	// 验证必填字段
	// Convert JSON strings to native types
	tags := unmarshalJSONStringSlice(in.Tags)
	parameters := unmarshalJSONMap(in.Parameters)
	targetCiTypes := unmarshalJSONStringSlice(in.TargetCiTypes)
	targetOsTypes := unmarshalJSONStringSlice(in.TargetOsTypes)
	targetSelector := unmarshalJSONMap(in.TargetSelector)
	requiredCapabilities := unmarshalJSONStringSlice(in.RequiredCapabilities)
	defaultEnv := unmarshalJSONStringMap(in.DefaultEnv)

	if in.Name == nil || *in.Name == "" {
		return nil, newInvalidArgumentError("script name is required")
	}
	if in.Code == nil || *in.Code == "" {
		return nil, newInvalidArgumentError("script code is required")
	}
	if in.ScriptType == nil || *in.ScriptType == "" {
		return nil, newInvalidArgumentError("script type is required")
	}
	if in.Executor == nil || *in.Executor == "" {
		return nil, newInvalidArgumentError("executor is required")
	}
	if in.Content == nil || *in.Content == "" {
		return nil, newInvalidArgumentError("script content is required")
	}

	// 设置默认值
	if in.Version == nil || *in.Version == "" {
		in.Version = pointy.GetPointer("1.0.0")
	}
	if in.IsLatest == nil {
		in.IsLatest = pointy.GetPointer(true)
	}
	if in.ExecutionCount == nil {
		in.ExecutionCount = pointy.GetPointer(uint64(0))
	}
	if in.SuccessCount == nil {
		in.SuccessCount = pointy.GetPointer(uint64(0))
	}
	if in.FailureCount == nil {
		in.FailureCount = pointy.GetPointer(uint64(0))
	}

	// 创建脚本
	query := l.svcCtx.DB.Script.Create().
		SetNotNilName(in.Name).
		SetNotNilCode(in.Code).
		SetNotNilDescription(in.Description).
		SetNotNilCategoryID(in.CategoryId).
		SetTags(tags).
		SetNotNilScriptType(in.ScriptType).
		SetNotNilExecutor(in.Executor).
		SetNotNilContent(in.Content).
		SetNotNilIsTemplate(in.IsTemplate).
		SetNotNilTemplateEngine(in.TemplateEngine).
		SetParameters(parameters).
		SetNotNilVersion(in.Version).
		SetNotNilBaseVersionID(in.BaseVersionId).
		SetNotNilIsLatest(in.IsLatest).
		SetTargetCiTypes(targetCiTypes).
		SetTargetOsTypes(targetOsTypes).
		SetTargetSelector(targetSelector).
		SetNotNilCredentialRef(in.CredentialRef).
		SetRequiredCapabilities(requiredCapabilities).
		SetNotNilDefaultTimeout(in.DefaultTimeout).
		SetNotNilDefaultWorkdir(in.DefaultWorkdir).
		SetDefaultEnv(defaultEnv).
		SetNotNilRequireConfirmation(in.RequireConfirmation).
		SetNotNilRiskLevel(in.RiskLevel).
		SetNotNilSchedulable(in.Schedulable).
		SetNotNilDefaultSchedule(in.DefaultSchedule).
		SetNotNilExecutionCount(in.ExecutionCount).
		SetNotNilSuccessCount(in.SuccessCount).
		SetNotNilFailureCount(in.FailureCount).
		SetNotNilLastExecutedAt(in.LastExecutedAt)

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}

	result, err := query.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 创建初始版本记录
	if in.Version != nil && *in.Version != "" {
		_, vErr := l.svcCtx.DB.ScriptVersion.Create().
			SetScriptID(result.ID).
			SetVersion(*in.Version).
			SetContent(*in.Content).
			SetParameters(parameters).
			SetNotNilScriptType(in.ScriptType).
			SetNotNilExecutor(in.Executor).
			SetChangeLog("Initial version").
			Save(l.ctx)
		if vErr != nil {
			l.Errorf("Failed to create script version: %v", vErr)
		}
	}

	return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess}, nil
}
