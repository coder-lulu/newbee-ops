package agent

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Agent
func NewDeleteAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentLogic {
	return &DeleteAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAgentLogic) DeleteAgent(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
