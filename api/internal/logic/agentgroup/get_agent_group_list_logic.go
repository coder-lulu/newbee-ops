package agentgroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Agent分组列表
func NewGetAgentGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentGroupListLogic {
	return &GetAgentGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentGroupListLogic) GetAgentGroupList(req *types.AgentGroupListReq) (resp *types.AgentGroupListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
