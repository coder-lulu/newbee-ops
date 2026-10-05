package agentgroup

import (
	"context"
	pb "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent分组详情
func NewGetAgentGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupByIdLogic {
	return &GetAgentGroupByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentGroupByIdLogic) GetAgentGroupById(req *types.IDReq) (resp *types.AgentGroupDetail, err error) {
	item, err := l.svcCtx.OpsClient.GetAgentGroupById(l.ctx, &pb.IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.AgentGroupDetail{Id: item.GetId(),
		CreatedAt:           item.GetCreatedAt(),
		UpdatedAt:           item.GetUpdatedAt(),
		Status:              item.GetStatus(),
		Name:                item.GetName(),
		Description:         item.GetDescription(),
		SelectionStrategy:   item.GetSelectionStrategy(),
		HealthCheckInterval: item.GetHealthCheckInterval(),
		AutoFailover:        item.GetAutoFailover(),
		MaxRetryCount:       item.GetMaxRetryCount()}, nil
}
