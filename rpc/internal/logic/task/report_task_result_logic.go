package task

import (
	"entgo.io/ent/dialect/sql"
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
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
	// Get task by TaskID (which is a string, not pointer)
	task, err := l.svcCtx.DB.Task.Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.EQ(s.C("task_id"), in.TaskId))
		}).
		First(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// Update task with result (TaskId and Status are strings, others are pointers)
	updateQuery := l.svcCtx.DB.Task.UpdateOneID(task.ID).
		SetStatusStr(in.Status)

	if in.ResultData != nil {
		updateQuery = updateQuery.SetResultData(*in.ResultData)
	}
	if in.ErrorMessage != nil {
		updateQuery = updateQuery.SetErrorMsg(*in.ErrorMessage)
	}
	if in.ExecutionTimeMs != nil {
		updateQuery = updateQuery.SetExecutionTimeMs(*in.ExecutionTimeMs)
	}

	_, err = updateQuery.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// TODO: ScriptExecutor temporarily disabled - Phase 3 feature
	// if l.svcCtx.ScriptExecutor != nil {
	// 	err = l.svcCtx.ScriptExecutor.UpdateTaskResult(
	// 		l.ctx,
	// 		in.TaskId,
	// 		in.Status,
	// 		in.ResultData,
	// 		in.ErrorMessage,
	// 		in.ExecutionTimeMs,
	// 	)
	// 	if err != nil {
	// 		l.Errorw("Failed to aggregate script execution result",
	// 			logx.Field("task_id", in.TaskId),
	// 			logx.Field("error", err))
	// 	}
	// }

	return &ops.BaseResp{Msg: "Task result reported successfully"}, nil
}
