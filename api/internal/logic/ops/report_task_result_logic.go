package ops

import (
	"context"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportTaskResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportTaskResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportTaskResultLogic {
	return &ReportTaskResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportTaskResultLogic) ReportTaskResult(req *types.ReportTaskResultReq) (resp *types.BaseMsgResp, err error) {
	l.Infow("Received task result report from Proxy",
		logx.Field("task_id", req.TaskId),
		logx.Field("status", req.Status))

	// 调用RPC上报任务结果
	result, err := l.svcCtx.OpsClient.ReportTaskResult(l.ctx, req.TaskId, req.Status, req.ResultData, req.ErrorMessage, req.ExecutionTimeMs, req.EndTime)
	if err != nil {
		l.Errorw("Failed to report task result to RPC",
			logx.Field("task_id", req.TaskId),
			logx.Field("error", err))
		return &types.BaseMsgResp{
			Code: 1,
			Msg:  "Failed to report task result: " + err.Error(),
		}, err
	}

	l.Infow("Task result reported successfully",
		logx.Field("task_id", req.TaskId),
		logx.Field("status", req.Status))

	return &types.BaseMsgResp{
		Code: 0,
		Msg:  result.Msg,
	}, nil
}
