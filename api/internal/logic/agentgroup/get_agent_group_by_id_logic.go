package agentgroup

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
