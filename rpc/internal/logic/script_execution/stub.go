package script_execution

import (
	"context"
	"errors"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

// Stub implementations - Phase 2/3 features temporarily disabled

type ExecuteScriptLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewExecuteScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteScriptLogic {
	return &ExecuteScriptLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *ExecuteScriptLogic) ExecuteScript(in *ops.ExecuteScriptReq) (*ops.ExecuteScriptResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetExecutionStatusLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetExecutionStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionStatusLogic {
	return &GetExecutionStatusLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetExecutionStatusLogic) GetExecutionStatus(in *ops.GetExecutionStatusReq) (*ops.ScriptExecutionInfo, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetExecutionDetailsLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetExecutionDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionDetailsLogic {
	return &GetExecutionDetailsLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetExecutionDetailsLogic) GetExecutionDetails(in *ops.GetExecutionDetailsReq) (*ops.GetExecutionDetailsResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type CancelExecutionLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewCancelExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelExecutionLogic {
	return &CancelExecutionLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *CancelExecutionLogic) CancelExecution(in *ops.CancelExecutionReq) (*ops.BaseResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type CreateScriptExecutionLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewCreateScriptExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptExecutionLogic {
	return &CreateScriptExecutionLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *CreateScriptExecutionLogic) CreateScriptExecution(in *ops.ScriptExecutionInfo) (*ops.BaseIDResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type UpdateScriptExecutionLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewUpdateScriptExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptExecutionLogic {
	return &UpdateScriptExecutionLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *UpdateScriptExecutionLogic) UpdateScriptExecution(in *ops.ScriptExecutionInfo) (*ops.BaseResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetScriptExecutionListLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetScriptExecutionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptExecutionListLogic {
	return &GetScriptExecutionListLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetScriptExecutionListLogic) GetScriptExecutionList(in *ops.ScriptExecutionListReq) (*ops.ScriptExecutionListResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetScriptExecutionByIdLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetScriptExecutionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptExecutionByIdLogic {
	return &GetScriptExecutionByIdLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetScriptExecutionByIdLogic) GetScriptExecutionById(in *ops.IDReq) (*ops.ScriptExecutionInfo, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type DeleteScriptExecutionLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewDeleteScriptExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScriptExecutionLogic {
	return &DeleteScriptExecutionLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *DeleteScriptExecutionLogic) DeleteScriptExecution(in *ops.IDsReq) (*ops.BaseResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}
