package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/ent/agentgroup"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/msg/errormsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentGroupLogic {
	return &CreateAgentGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentGroupLogic) CreateAgentGroup(in *ops.AgentGroupInfo) (*ops.BaseIDResp, error) {
	query := l.svcCtx.DB.AgentGroup.Create().
		SetNotNilName(in.Name).
		SetNotNilDescription(in.Description).
		SetNotNilAutoFailover(in.AutoFailover)

	// 转换enum字段
	if in.SelectionStrategy != nil {
		query.SetSelectionStrategy(agentgroup.SelectionStrategy(*in.SelectionStrategy))
	}

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.HealthCheckInterval != nil {
		query.SetNotNilHealthCheckInterval(pointy.GetPointer(int(*in.HealthCheckInterval)))
	}
	if in.MaxRetryCount != nil {
		query.SetNotNilMaxRetryCount(pointy.GetPointer(int(*in.MaxRetryCount)))
	}

	result, err := query.Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &ops.BaseIDResp{Id: result.ID, Msg: errormsg.CreateSuccess}, nil
}
