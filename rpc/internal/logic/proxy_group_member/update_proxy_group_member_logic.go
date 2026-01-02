package proxy_group_member

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/suyuan32/simple-admin-common/msg/errormsg"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProxyGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProxyGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProxyGroupMemberLogic {
	return &UpdateProxyGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProxyGroupMemberLogic) UpdateProxyGroupMember(in *ops.ProxyGroupMemberInfo) (*ops.BaseResp, error) {
	query:= l.svcCtx.DB.ProxyGroupMember.UpdateOneID(*in.Id).
			SetNotNilDepartmentID(in.DepartmentId).
			SetNotNilWorkerID(in.ProxyId).
			SetNotNilWorkerGroupID(in.ProxyGroupId).
			SetNotNilJoinedAt(pointy.GetTimeMilliPointer(in.JoinedAt)).
			SetNotNilLastSelectedAt(pointy.GetTimeMilliPointer(in.LastSelectedAt)).
			SetNotNilIsBackup(in.IsBackup).
			SetNotNilLastFailureAt(pointy.GetTimeMilliPointer(in.LastFailureAt))

	if in.Status != nil {
		query.SetNotNilStatus(pointy.GetPointer(uint8(*in.Status)))
	}
	if in.Weight != nil {
		query.SetNotNilWeight(pointy.GetPointer(int(*in.Weight)))
	}
	if in.Priority != nil {
		query.SetNotNilPriority(pointy.GetPointer(int(*in.Priority)))
	}
	if in.SelectCount != nil {
		query.SetNotNilSelectCount(pointy.GetPointer(int(*in.SelectCount)))
	}
	if in.ConsecutiveFailures != nil {
		query.SetNotNilConsecutiveFailures(pointy.GetPointer(int(*in.ConsecutiveFailures)))
	}

	 err := query.Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &ops.BaseResp{Msg: errormsg.UpdateSuccess }, nil
}
