package agent

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent详情
func NewGetAgentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentByIdLogic {
	return &GetAgentByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentByIdLogic) GetAgentById(req *types.IDReq) (resp *types.AgentDetail, err error) {
	// todo: add your logic here and delete this line

	return
}
