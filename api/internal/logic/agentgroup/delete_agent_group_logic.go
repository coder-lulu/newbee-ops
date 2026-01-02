package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Agent分组
func NewDeleteAgentGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentGroupLogic {
	return &DeleteAgentGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAgentGroupLogic) DeleteAgentGroup(req *types.IDsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
