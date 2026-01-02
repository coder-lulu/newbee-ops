package workergroup

import (
	"context"
	"errors"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
)

// Stub implementations - WorkerGroup features need manual update to ProxyGroup

var errNotImplemented = errors.New("WorkerGroup feature temporarily disabled - needs manual migration to ProxyGroup")

type CreateWorkerGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewCreateWorkerGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkerGroupLogic {
	return &CreateWorkerGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *CreateWorkerGroupLogic) CreateWorkerGroup(in *ops.WorkerGroupInfo) (*ops.BaseIDResp, error) {
	return nil, errNotImplemented
}

type UpdateWorkerGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewUpdateWorkerGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkerGroupLogic {
	return &UpdateWorkerGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *UpdateWorkerGroupLogic) UpdateWorkerGroup(in *ops.WorkerGroupInfo) (*ops.BaseResp, error) {
	return nil, errNotImplemented
}

type GetWorkerGroupListLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetWorkerGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkerGroupListLogic {
	return &GetWorkerGroupListLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetWorkerGroupListLogic) GetWorkerGroupList(in *ops.WorkerGroupListReq) (*ops.WorkerGroupListResp, error) {
	return nil, errNotImplemented
}

type GetWorkerGroupByIdLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetWorkerGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkerGroupByIdLogic {
	return &GetWorkerGroupByIdLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetWorkerGroupByIdLogic) GetWorkerGroupById(in *ops.IDReq) (*ops.WorkerGroupInfo, error) {
	return nil, errNotImplemented
}

type DeleteWorkerGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewDeleteWorkerGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkerGroupLogic {
	return &DeleteWorkerGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *DeleteWorkerGroupLogic) DeleteWorkerGroup(in *ops.IDsReq) (*ops.BaseResp, error) {
	return nil, errNotImplemented
}

type AddWorkersToGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewAddWorkersToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWorkersToGroupLogic {
	return &AddWorkersToGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *AddWorkersToGroupLogic) AddWorkersToGroup(in *ops.WorkerGroupMemberReq) (*ops.BaseResp, error) {
	return nil, errNotImplemented
}

type RemoveWorkersFromGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewRemoveWorkersFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveWorkersFromGroupLogic {
	return &RemoveWorkersFromGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *RemoveWorkersFromGroupLogic) RemoveWorkersFromGroup(in *ops.WorkerGroupMemberReq) (*ops.BaseResp, error) {
	return nil, errNotImplemented
}

type UpdateWorkerGroupMemberWeightLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewUpdateWorkerGroupMemberWeightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkerGroupMemberWeightLogic {
	return &UpdateWorkerGroupMemberWeightLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *UpdateWorkerGroupMemberWeightLogic) UpdateWorkerGroupMemberWeight(in *ops.WorkerGroupMemberInfo) (*ops.BaseResp, error) {
	return nil, errNotImplemented
}

type GetWorkerGroupMembersLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewGetWorkerGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkerGroupMembersLogic {
	return &GetWorkerGroupMembersLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *GetWorkerGroupMembersLogic) GetWorkerGroupMembers(in *ops.IDReq) (*ops.WorkerGroupInfoWithMembers, error) {
	return nil, errNotImplemented
}

type PickWorkerFromGroupLogic struct{ ctx context.Context; svcCtx *svc.ServiceContext }
func NewPickWorkerFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PickWorkerFromGroupLogic {
	return &PickWorkerFromGroupLogic{ctx: ctx, svcCtx: svcCtx}
}
func (l *PickWorkerFromGroupLogic) PickWorkerFromGroup(in *ops.PickWorkerFromGroupReq) (*ops.PickWorkerFromGroupResp, error) {
	return nil, errNotImplemented
}
