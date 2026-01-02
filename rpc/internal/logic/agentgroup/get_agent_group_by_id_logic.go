package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentGroupByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupByIdLogic {
	return &GetAgentGroupByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentGroupByIdLogic) GetAgentGroupById(in *ops.IDReq) (*ops.AgentGroupInfo, error) {
	result, err := l.svcCtx.DB.AgentGroup.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 转换enum为string
	selectionStrategyStr := string(result.SelectionStrategy)

	return &ops.AgentGroupInfo{
		Id:                  &result.ID,
		CreatedAt:           pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:           pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:              pointy.GetPointer(uint32(result.Status)),
		Name:                &result.Name,
		Description:         &result.Description,
		SelectionStrategy:   &selectionStrategyStr,
		HealthCheckInterval: pointy.GetPointer(int64(result.HealthCheckInterval)),
		AutoFailover:        &result.AutoFailover,
		MaxRetryCount:       pointy.GetPointer(int64(result.MaxRetryCount)),
	}, nil
}

