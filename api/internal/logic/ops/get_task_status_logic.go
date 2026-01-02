package ops

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskStatusLogic {
	return &GetTaskStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskStatusLogic) GetTaskStatus(req *types.GetTaskStatusReq) (resp *types.TaskStatusResp, err error) {
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
		// 如果StatusStr为空，尝试从Status数字转换
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

	l.Logger.Infow("Got task status",
		logx.Field("task_id", req.TaskId),
		logx.Field("status", status))

	return &types.TaskStatusResp{
		Code: 0,
		Msg:  "success",
		Data: types.TaskStatusData{
			TaskId: req.TaskId,
			Status: status,
		},
	}, nil
}
