package ops

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskResultLogic {
	return &GetTaskResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskResultLogic) GetTaskResult(req *types.GetTaskResultReq) (resp *types.TaskResultResp, err error) {
	// 1. 查询任务信息
	taskInfo, err := l.svcCtx.OpsClient.GetTaskByTaskId(l.ctx, &ops.TaskIdReq{
		TaskId: req.TaskId,
	})
	if err != nil {
		l.Logger.Errorw("Failed to get task by task_id",
			logx.Field("task_id", req.TaskId),
			logx.Field("error", err))
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// 2. 转换状态
	status := "unknown"
	if taskInfo.StatusStr != nil {
		status = *taskInfo.StatusStr
	} else if taskInfo.Status != nil {
		switch *taskInfo.Status {
		case 0:
			status = "pending"
		case 1:
			status = "running"
		case 2:
			status = "completed"
		case 3:
			status = "failed"
		case 4:
			status = "cancelled"
		default:
			status = "unknown"
		}
	}

	// 3. 提取结果信息
	output := ""
	if taskInfo.ResultOutput != nil {
		output = *taskInfo.ResultOutput
	}

	errorMsg := ""
	if taskInfo.ErrorMsg != nil {
		errorMsg = *taskInfo.ErrorMsg
	}

	exitCode := int(0)
	if taskInfo.ExitCode != nil {
		exitCode = int(*taskInfo.ExitCode)
	}

	l.Logger.Infow("Got task result",
		logx.Field("task_id", req.TaskId),
		logx.Field("status", status),
		logx.Field("exit_code", exitCode),
		logx.Field("has_output", len(output) > 0),
		logx.Field("has_error", len(errorMsg) > 0))

	return &types.TaskResultResp{
		Code: 0,
		Msg:  "success",
		Data: types.TaskResultData{
			TaskId:   req.TaskId,
			Status:   status,
			Output:   output,
			Error:    errorMsg,
			ExitCode: exitCode,
		},
	}, nil
}
