package agentgroupmember

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroupmember"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupMemberListLogic {
	return &GetAgentGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentGroupMemberListLogic) GetAgentGroupMemberList(in *ops.AgentGroupMemberListReq) (*ops.AgentGroupMemberListResp, error) {
	var predicates []predicate.AgentGroupMember
	if in.CreatedAt != nil {
		predicates = append(predicates, agentgroupmember.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, agentgroupmember.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, agentgroupmember.StatusEQ(uint8(*in.Status)))
	}
	if in.AgentId != nil {
		predicates = append(predicates, agentgroupmember.AgentIDEQ(*in.AgentId))
	}
	if in.AgentGroupId != nil {
		predicates = append(predicates, agentgroupmember.AgentGroupIDEQ(*in.AgentGroupId))
	}
	if in.Priority != nil {
		predicates = append(predicates, agentgroupmember.PriorityEQ(int(*in.Priority)))
	}
	if in.Weight != nil {
		predicates = append(predicates, agentgroupmember.WeightEQ(int(*in.Weight)))
	}
	if in.JoinedAt != nil {
		predicates = append(predicates, agentgroupmember.JoinedAtGTE(time.UnixMilli(*in.JoinedAt)))
	}
	result, err := l.svcCtx.DB.AgentGroupMember.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.AgentGroupMemberListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.AgentGroupMemberInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:	pointy.GetPointer(uint32(v.Status)),
			AgentId:	&v.AgentID,
			AgentGroupId:	&v.AgentGroupID,
			Priority:	pointy.GetPointer(int64(v.Priority)),
			Weight:	pointy.GetPointer(int64(v.Weight)),
			JoinedAt:	pointy.GetUnixMilliPointer(v.JoinedAt.UnixMilli()),
		})
	}

	return resp, nil
}
