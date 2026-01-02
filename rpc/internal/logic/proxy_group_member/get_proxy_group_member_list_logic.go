package proxy_group_member

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/proxygroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupMemberListLogic {
	return &GetProxyGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyGroupMemberListLogic) GetProxyGroupMemberList(in *ops.ProxyGroupMemberListReq) (*ops.ProxyGroupMemberListResp, error) {
	var predicates []predicate.ProxyGroupMember
	if in.CreatedAt != nil {
		predicates = append(predicates, proxygroupmember.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, proxygroupmember.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.DepartmentId != nil {
		predicates = append(predicates, proxygroupmember.DepartmentIDEQ(*in.DepartmentId))
	}
	if in.Status != nil {
		predicates = append(predicates, proxygroupmember.StatusEQ(uint8(*in.Status)))
	}
	if in.ProxyId != nil {
		predicates = append(predicates, proxygroupmember.WorkerIDEQ(*in.ProxyId))
	}
	if in.ProxyGroupId != nil {
		predicates = append(predicates, proxygroupmember.WorkerGroupIDEQ(*in.ProxyGroupId))
	}
	if in.Weight != nil {
		predicates = append(predicates, proxygroupmember.WeightEQ(int(*in.Weight)))
	}
	if in.Priority != nil {
		predicates = append(predicates, proxygroupmember.PriorityEQ(int(*in.Priority)))
	}
	if in.JoinedAt != nil {
		predicates = append(predicates, proxygroupmember.JoinedAtGTE(time.UnixMilli(*in.JoinedAt)))
	}
	if in.LastSelectedAt != nil {
		predicates = append(predicates, proxygroupmember.LastSelectedAtGTE(time.UnixMilli(*in.LastSelectedAt)))
	}
	if in.SelectCount != nil {
		predicates = append(predicates, proxygroupmember.SelectCountEQ(int(*in.SelectCount)))
	}
	if in.IsBackup != nil {
		predicates = append(predicates, proxygroupmember.IsBackupEQ(*in.IsBackup))
	}
	if in.ConsecutiveFailures != nil {
		predicates = append(predicates, proxygroupmember.ConsecutiveFailuresEQ(int(*in.ConsecutiveFailures)))
	}
	if in.LastFailureAt != nil {
		predicates = append(predicates, proxygroupmember.LastFailureAtGTE(time.UnixMilli(*in.LastFailureAt)))
	}
	result, err := l.svcCtx.DB.ProxyGroupMember.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ProxyGroupMemberListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.ProxyGroupMemberInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			DepartmentId:	&v.DepartmentID,
			Status:	pointy.GetPointer(uint32(v.Status)),
			ProxyId:	&v.WorkerID,
			ProxyGroupId:	&v.WorkerGroupID,
			Weight:	pointy.GetPointer(int64(v.Weight)),
			Priority:	pointy.GetPointer(int64(v.Priority)),
			JoinedAt:	pointy.GetUnixMilliPointer(v.JoinedAt.UnixMilli()),
			LastSelectedAt:	pointy.GetUnixMilliPointer(v.LastSelectedAt.UnixMilli()),
			SelectCount:	pointy.GetPointer(int64(v.SelectCount)),
			IsBackup:	&v.IsBackup,
			ConsecutiveFailures:	pointy.GetPointer(int64(v.ConsecutiveFailures)),
			LastFailureAt:	pointy.GetUnixMilliPointer(v.LastFailureAt.UnixMilli()),
		})
	}

	return resp, nil
}
