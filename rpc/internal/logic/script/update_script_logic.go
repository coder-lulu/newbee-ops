package script

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/msg/errormsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScriptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptLogic {
	return &UpdateScriptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScriptLogic) UpdateScript(in *ops.ScriptInfo) (*ops.BaseResp, error) {
	// Convert JSON strings to native types
	tags := unmarshalJSONStringSlice(in.Tags)
	parameters := unmarshalJSONMap(in.Parameters)
	targetCiTypes := unmarshalJSONStringSlice(in.TargetCiTypes)
	targetOsTypes := unmarshalJSONStringSlice(in.TargetOsTypes)
	targetSelector := unmarshalJSONMap(in.TargetSelector)
	requiredCapabilities := unmarshalJSONStringSlice(in.RequiredCapabilities)
	defaultEnv := unmarshalJSONStringMap(in.DefaultEnv)
	if in.Id == nil {
		return nil, newInvalidArgumentError("script id is required")
	}

	// 获取当前脚本信息，用于版本对比
	oldScript, err := l.svcCtx.DB.Script.Get(l.ctx, *in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 检查内容是否变更，如果变更则创建新版本
	contentChanged := false
	if in.Content != nil && *in.Content != oldScript.Content {
		contentChanged = true
	}

	query := l.svcCtx.DB.Script.UpdateOneID(*in.Id).
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

	// 如果内容变更，创建新版本记录
	if contentChanged && in.Version != nil && *in.Version != "" {
		_, vErr := l.svcCtx.DB.ScriptVersion.Create().
			SetScriptID(result.ID).
			SetVersion(*in.Version).
			SetContent(*in.Content).
			SetParameters(parameters).
			SetNotNilScriptType(in.ScriptType).
			SetNotNilExecutor(in.Executor).
			SetChangeLog("Updated at " + time.Now().Format("2006-01-02 15:04:05")).
			Save(l.ctx)
		if vErr != nil {
			l.Errorf("Failed to create script version: %v", vErr)
		}
	}

	return &ops.BaseResp{Msg: errormsg.UpdateSuccess}, nil
}
