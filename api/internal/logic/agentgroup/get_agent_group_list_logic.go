package agentgroup

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent分组列表
func NewGetAgentGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupListLogic {
	return &GetAgentGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentGroupListLogic) GetAgentGroupList(req *types.AgentGroupListReq) (resp *types.AgentGroupListResp, err error) {
	rpcReq := &pb.AgentGroupListReq{Page: req.Page, PageSize: req.PageSize}
	if req.Name != "" {
		rpcReq.Name = &req.Name
	}
	if req.Status != 0 {
		rpcReq.Status = &req.Status
	}
	result, err := l.svcCtx.OpsClient.GetAgentGroupList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	items := make([]types.AgentGroupItem, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, types.AgentGroupItem{Id: item.GetId(),
			Name:                item.GetName(),
			Description:         item.GetDescription(),
			SelectionStrategy:   item.GetSelectionStrategy(),
			HealthCheckInterval: item.GetHealthCheckInterval(),
			AutoFailover:        item.GetAutoFailover()})
	}
	return &types.AgentGroupListResp{Total: result.Total, Data: items}, nil
}
