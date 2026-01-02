package agent

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectAgentForDiscoveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 选择Agent（供Discovery使用）
func NewSelectAgentForDiscoveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAgentForDiscoveryLogic {
	return &SelectAgentForDiscoveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectAgentForDiscoveryLogic) SelectAgentForDiscovery(req *types.AgentSelectionReq) (resp *types.AgentSelectionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
