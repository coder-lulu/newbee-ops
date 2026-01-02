package proxygroup

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新分组成员（需要JWT）- 替换式更新，自动处理增删
func NewUpdateGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupMembersLogic {
	return &UpdateGroupMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGroupMembersLogic) UpdateGroupMembers(req *types.UpdateGroupMembersReq) (resp *types.UpdateGroupMembersResp, err error) {
	// todo: add your logic here and delete this line

	return
}
