package script

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetScriptByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptByIdLogic {
	return &GetScriptByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptByIdLogic) GetScriptById(in *ops.IDReq) (*ops.ScriptInfo, error) {
	result, err := l.svcCtx.DB.Script.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.ScriptInfo{
		Id:                   &result.ID,
		CreatedAt:            pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:            pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:               pointy.GetPointer(uint32(result.Status)),
		Name:                 &result.Name,
		Code:                 &result.Code,
		Description:          &result.Description,
		CategoryId:           &result.CategoryID,
		Tags: pointy.GetPointer(mustMarshalJSON(result.Tags)),
		ScriptType:           &result.ScriptType,
		Executor:             &result.Executor,
		Content:              &result.Content,
		IsTemplate:           &result.IsTemplate,
		TemplateEngine:       &result.TemplateEngine,
		Parameters: pointy.GetPointer(mustMarshalJSON(result.Parameters)),
		Version:              &result.Version,
		BaseVersionId:        &result.BaseVersionID,
		IsLatest:             &result.IsLatest,
		TargetCiTypes: pointy.GetPointer(mustMarshalJSON(result.TargetCiTypes)),
		TargetOsTypes: pointy.GetPointer(mustMarshalJSON(result.TargetOsTypes)),
		TargetSelector: pointy.GetPointer(mustMarshalJSON(result.TargetSelector)),
		CredentialRef:        &result.CredentialRef,
		RequiredCapabilities: pointy.GetPointer(mustMarshalJSON(result.RequiredCapabilities)),
		DefaultTimeout:       &result.DefaultTimeout,
		DefaultWorkdir:       &result.DefaultWorkdir,
		DefaultEnv: pointy.GetPointer(mustMarshalJSON(result.DefaultEnv)),
		RequireConfirmation:  &result.RequireConfirmation,
		RiskLevel:            &result.RiskLevel,
		Schedulable:          &result.Schedulable,
		DefaultSchedule:      &result.DefaultSchedule,
		ExecutionCount:       &result.ExecutionCount,
		SuccessCount:         &result.SuccessCount,
		FailureCount:         &result.FailureCount,
		LastExecutedAt:       &result.LastExecutedAt,
	}, nil
}
