package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建Agent分组
func NewCreateAgentGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentGroupLogic {
	return &CreateAgentGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAgentGroupLogic) CreateAgentGroup(req *types.AgentGroupInfo) (resp *types.BaseIDResp, err error) {
	// todo: add your logic here and delete this line

	return
}
