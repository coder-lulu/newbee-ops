package agentgroup

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroup"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupListLogic {
	return &GetAgentGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentGroupListLogic) GetAgentGroupList(in *ops.AgentGroupListReq) (*ops.AgentGroupListResp, error) {
	var predicates []predicate.AgentGroup
	if in.CreatedAt != nil {
		predicates = append(predicates, agentgroup.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, agentgroup.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, agentgroup.StatusEQ(uint8(*in.Status)))
	}
	if in.Name != nil {
		predicates = append(predicates, agentgroup.NameContains(*in.Name))
	}
	if in.Description != nil {
		predicates = append(predicates, agentgroup.DescriptionContains(*in.Description))
	}
	if in.SelectionStrategy != nil {
		predicates = append(predicates, agentgroup.SelectionStrategyEQ(agentgroup.SelectionStrategy(*in.SelectionStrategy)))
	}
	if in.HealthCheckInterval != nil {
		predicates = append(predicates, agentgroup.HealthCheckIntervalEQ(int(*in.HealthCheckInterval)))
	}
	if in.AutoFailover != nil {
		predicates = append(predicates, agentgroup.AutoFailoverEQ(*in.AutoFailover))
	}
	if in.MaxRetryCount != nil {
		predicates = append(predicates, agentgroup.MaxRetryCountEQ(int(*in.MaxRetryCount)))
	}
	result, err := l.svcCtx.DB.AgentGroup.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.AgentGroupListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		// 转换enum为string
		selectionStrategyStr := string(v.SelectionStrategy)

		resp.Data = append(resp.Data, &ops.AgentGroupInfo{
			Id:                  &v.ID,
			CreatedAt:           pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:           pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:              pointy.GetPointer(uint32(v.Status)),
			Name:                &v.Name,
			Description:         &v.Description,
			SelectionStrategy:   &selectionStrategyStr,
			HealthCheckInterval: pointy.GetPointer(int64(v.HealthCheckInterval)),
			AutoFailover:        &v.AutoFailover,
			MaxRetryCount:       pointy.GetPointer(int64(v.MaxRetryCount)),
		})
	}

	return resp, nil
}
