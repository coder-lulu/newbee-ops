package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupMemberWeightLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分组成员权重（需要JWT）
func NewUpdateGroupMemberWeightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupMemberWeightLogic {
	return &UpdateGroupMemberWeightLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGroupMemberWeightLogic) UpdateGroupMemberWeight(req *types.UpdateGroupMemberWeightReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
