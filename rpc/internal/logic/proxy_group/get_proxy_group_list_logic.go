package proxy_group

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/proxygroup"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetProxyGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProxyGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProxyGroupListLogic {
	return &GetProxyGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProxyGroupListLogic) GetProxyGroupList(in *ops.ProxyGroupListReq) (*ops.ProxyGroupListResp, error) {
	var predicates []predicate.ProxyGroup
	if in.CreatedAt != nil {
		predicates = append(predicates, proxygroup.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, proxygroup.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.DepartmentId != nil {
		predicates = append(predicates, proxygroup.DepartmentIDEQ(*in.DepartmentId))
	}
	if in.Status != nil {
		predicates = append(predicates, proxygroup.StatusEQ(uint8(*in.Status)))
	}
	if in.Name != nil {
		predicates = append(predicates, proxygroup.NameContains(*in.Name))
	}
	if in.Description != nil {
		predicates = append(predicates, proxygroup.DescriptionContains(*in.Description))
	}
	if in.SelectionStrategy != nil {
		predicates = append(predicates, proxygroup.SelectionStrategyEQ(stringToSelectionStrategy(in.SelectionStrategy)))
	}
	if in.HealthCheckInterval != nil {
		predicates = append(predicates, proxygroup.HealthCheckIntervalEQ(int(*in.HealthCheckInterval)))
	}
	if in.AutoFailover != nil {
		predicates = append(predicates, proxygroup.AutoFailoverEQ(*in.AutoFailover))
	}
	if in.MaxRetryCount != nil {
		predicates = append(predicates, proxygroup.MaxRetryCountEQ(int(*in.MaxRetryCount)))
	}
	if in.MinHealthyProxies != nil {
		predicates = append(predicates, proxygroup.MinHealthyWorkersEQ(int(*in.MinHealthyProxies)))
	}
	if in.ConnectionTimeout != nil {
		predicates = append(predicates, proxygroup.ConnectionTimeoutEQ(int(*in.ConnectionTimeout)))
	}
	if in.RequestTimeout != nil {
		predicates = append(predicates, proxygroup.RequestTimeoutEQ(int(*in.RequestTimeout)))
	}
	if in.TotalMembers != nil {
		predicates = append(predicates, proxygroup.TotalMembersEQ(int(*in.TotalMembers)))
	}
	if in.OnlineMembers != nil {
		predicates = append(predicates, proxygroup.OnlineMembersEQ(int(*in.OnlineMembers)))
	}
	if in.TotalWeight != nil {
		predicates = append(predicates, proxygroup.TotalWeightEQ(int(*in.TotalWeight)))
	}
	result, err := l.svcCtx.DB.ProxyGroup.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ProxyGroupListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.ProxyGroupInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			DepartmentId:	&v.DepartmentID,
			Status:	pointy.GetPointer(uint32(v.Status)),
			Name:	&v.Name,
			Description:	&v.Description,
			SelectionStrategy: pointy.GetPointer(selectionStrategyToString(v.SelectionStrategy)),
			HealthCheckInterval:	pointy.GetPointer(int64(v.HealthCheckInterval)),
			AutoFailover:	&v.AutoFailover,
			MaxRetryCount:	pointy.GetPointer(int64(v.MaxRetryCount)),
			MinHealthyProxies:	pointy.GetPointer(int64(v.MinHealthyWorkers)),
			ConnectionTimeout:	pointy.GetPointer(int64(v.ConnectionTimeout)),
			RequestTimeout:	pointy.GetPointer(int64(v.RequestTimeout)),
			TotalMembers:	pointy.GetPointer(int64(v.TotalMembers)),
			OnlineMembers:	pointy.GetPointer(int64(v.OnlineMembers)),
			TotalWeight:	pointy.GetPointer(int64(v.TotalWeight)),
		})
	}

	return resp, nil
}
