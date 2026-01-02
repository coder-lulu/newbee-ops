package accessprofile

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportTaskResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportTaskResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportTaskResultLogic {
	return &ReportTaskResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportTaskResultLogic) ReportTaskResult(in *ops.TaskResultReq) (*ops.BaseResp, error) {
	// todo: add your logic here and delete this line

	return &ops.BaseResp{}, nil
}
