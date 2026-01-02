package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新Agent分组
func NewUpdateAgentGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentGroupLogic {
	return &UpdateAgentGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAgentGroupLogic) UpdateAgentGroup(req *types.AgentGroupInfo) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
