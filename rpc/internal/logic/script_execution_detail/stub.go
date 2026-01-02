package script_execution_detail

import (
	"context"
	"errors"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

// Stub implementations - Phase 2/3 features temporarily disabled

type CreateScriptExecutionDetailLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewCreateScriptExecutionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScriptExecutionDetailLogic {
	return &CreateScriptExecutionDetailLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *CreateScriptExecutionDetailLogic) CreateScriptExecutionDetail(in *ops.ScriptExecutionDetailInfo) (*ops.BaseIDResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type UpdateScriptExecutionDetailLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewUpdateScriptExecutionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScriptExecutionDetailLogic {
	return &UpdateScriptExecutionDetailLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *UpdateScriptExecutionDetailLogic) UpdateScriptExecutionDetail(in *ops.ScriptExecutionDetailInfo) (*ops.BaseResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetScriptExecutionDetailListLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetScriptExecutionDetailListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptExecutionDetailListLogic {
	return &GetScriptExecutionDetailListLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetScriptExecutionDetailListLogic) GetScriptExecutionDetailList(in *ops.ScriptExecutionDetailListReq) (*ops.ScriptExecutionDetailListResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type GetScriptExecutionDetailByIdLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetScriptExecutionDetailByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptExecutionDetailByIdLogic {
	return &GetScriptExecutionDetailByIdLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetScriptExecutionDetailByIdLogic) GetScriptExecutionDetailById(in *ops.IDReq) (*ops.ScriptExecutionDetailInfo, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}

type DeleteScriptExecutionDetailLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewDeleteScriptExecutionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScriptExecutionDetailLogic {
	return &DeleteScriptExecutionDetailLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *DeleteScriptExecutionDetailLogic) DeleteScriptExecutionDetail(in *ops.IDsReq) (*ops.BaseResp, error) {
	return nil, errors.New("script execution feature temporarily disabled - pending Proxy architecture update")
}
