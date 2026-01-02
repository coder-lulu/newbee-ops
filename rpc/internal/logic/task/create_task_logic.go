package task

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/internal/logic/worker"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTaskLogic) CreateTask(in *newbee_ops_rpc.TaskInfo) (*newbee_ops_rpc.BaseIDResp, error) {
	q := l.svcCtx.DB.Task.Create()
	if in.TaskId != nil {
		q = q.SetTaskID(in.GetTaskId())
	}
	if in.Name != nil {
		q = q.SetName(in.GetName())
	}
	if in.CreatorId != nil {
		q = q.SetCreatorID(in.GetCreatorId())
	}
	// 转换structpb.Struct到[]string (ci_ids)
	if in.CiIds != nil {
		ciIdsAny := in.GetCiIds().AsMap()
		ciIds := make([]string, 0)
		if list, ok := ciIdsAny["list"].([]interface{}); ok {
			for _, v := range list {
				if str, ok := v.(string); ok {
					ciIds = append(ciIds, str)
				} else {
					l.Infow("skipping non-string ci_id value",
						logx.Field("type", fmt.Sprintf("%T", v)),
						logx.Field("value", v))
				}
			}
		}
		if len(ciIds) > 0 {
			q = q.SetCiIds(ciIds)
		}
	}
	if in.Executor != nil {
		q = q.SetExecutor(in.GetExecutor())
	}
	if in.CommandContent != nil {
		q = q.SetCommandContent(in.GetCommandContent())
	}
	if in.CommandTimeout != nil {
		timeout := in.GetCommandTimeout()
		// 验证int64→int转换不会溢出
		if timeout < math.MinInt32 || timeout > math.MaxInt32 {
			l.Errorw("command_timeout out of int32 range",
				logx.Field("timeout", timeout))
			return nil, fmt.Errorf("command_timeout out of range: %d (must be within %d to %d)", timeout, math.MinInt32, math.MaxInt32)
		}
		q = q.SetCommandTimeout(int(timeout))
	}
	if in.CommandWorkdir != nil {
		q = q.SetCommandWorkdir(in.GetCommandWorkdir())
	}
	// 转换structpb.Struct到map[string]string (command_env)
	if in.CommandEnv != nil {
		envAny := in.GetCommandEnv().AsMap()
		env := make(map[string]string)
		for k, v := range envAny {
			if str, ok := v.(string); ok {
				env[k] = str
			} else {
				l.Infow("skipping non-string command_env value",
					logx.Field("key", k),
					logx.Field("type", fmt.Sprintf("%T", v)),
					logx.Field("value", v))
			}
		}
		if len(env) > 0 {
			q = q.SetCommandEnv(env)
		}
	}
	if in.StatusStr != nil {
		q = q.SetStatusStr(in.GetStatusStr())
	}
	if in.ResultOutput != nil {
		q = q.SetResultOutput(in.GetResultOutput())
	}
	if in.ErrorMsg != nil {
		q = q.SetErrorMsg(in.GetErrorMsg())
	}
	if in.ExitCode != nil {
		exitCode := in.GetExitCode()
		// 验证int64→int转换不会溢出（exit code通常是0-255，但允许更大范围）
		if exitCode < math.MinInt32 || exitCode > math.MaxInt32 {
			l.Errorw("exit_code out of int32 range",
				logx.Field("exit_code", exitCode))
			return nil, fmt.Errorf("exit_code out of range: %d (must be within %d to %d)", exitCode, math.MinInt32, math.MaxInt32)
		}
		q = q.SetExitCode(int(exitCode))
	}
	if in.RemoteTaskId != nil {
		q = q.SetRemoteTaskID(in.GetRemoteTaskId())
	}
	if in.StartTime != nil {
		q = q.SetStartTime(in.GetStartTime())
	}
	if in.EndTime != nil {
		q = q.SetEndTime(in.GetEndTime())
	}
	// 转换structpb.Struct到map[string]string (tags)
	if in.Tags != nil {
		tagsAny := in.GetTags().AsMap()
		tags := make(map[string]string)
		for k, v := range tagsAny {
			if str, ok := v.(string); ok {
				tags[k] = str
			} else {
				l.Infow("skipping non-string tag value",
					logx.Field("key", k),
					logx.Field("type", fmt.Sprintf("%T", v)),
					logx.Field("value", v))
			}
		}
		if len(tags) > 0 {
			q = q.SetTags(tags)
		}
	}
	if in.Status != nil {
		q = q.SetStatus(uint8(in.GetStatus()))
	}

	// Phase 2: Worker选择和任务分配
	// 1. 选择Worker (如果executor需要Proxy执行)
	var selectedWorker *worker.WorkerInfo
	needsProxy := in.Executor != nil && (in.GetExecutor() == "ssh" || in.GetExecutor() == "telnet" || in.GetExecutor() == "rdp")

	if needsProxy {
		// 确定所需能力
		requiredCapabilities := []string{in.GetExecutor()}

		// 选择Worker - 使用least_connections策略
		workerInfo, err := l.svcCtx.WorkerManager.SelectWorker(l.ctx, 1, "least_connections", &worker.SelectOptions{
			RequiredCapabilities: requiredCapabilities,
		})
		if err != nil {
			l.Errorw("Failed to select worker for task",
				logx.Field("executor", in.GetExecutor()),
				logx.Field("error", err))
			// 不阻塞任务创建，标记为pending等待重试
		} else {
			selectedWorker = workerInfo
			q = q.SetWorkerID(workerInfo.WorkerID)
			l.Infow("Selected worker for task",
				logx.Field("worker_id", workerInfo.WorkerID),
				logx.Field("executor", in.GetExecutor()))
		}
	}

	// 设置初始状态
	if selectedWorker != nil {
		q = q.SetStatusStr("selecting_worker")
	} else if needsProxy {
		q = q.SetStatusStr("pending") // 等待Worker可用
	}

	// 2. 保存任务
	res, err := q.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 3. 异步分配任务到Proxy（如果已选择Worker）
	if selectedWorker != nil {
		go func() {
			dispatchCtx := context.Background()

			// 更新状态为dispatched
			_, err := l.svcCtx.DB.Task.UpdateOneID(res.ID).
				SetStatusStr("dispatching").
				SetDispatchedAt(time.Now()).
				Save(dispatchCtx)
			if err != nil {
				l.Errorw("Failed to update task status to dispatching",
					logx.Field("task_id", res.TaskID),
					logx.Field("error", err))
				return
			}

			// 调用TaskDispatcher分配任务
			// TaskDispatcher temporarily disabled
		err = fmt.Errorf("TaskDispatcher not available - needs Proxy architecture update") // l.svcCtx.TaskDispatcher.DispatchTask(dispatchCtx, res, selectedWorker)
			if err != nil {
				// 分配失败，更新状态
				l.Errorw("Failed to dispatch task to proxy",
					logx.Field("task_id", res.TaskID),
					logx.Field("worker_id", selectedWorker.WorkerID),
					logx.Field("error", err))

				l.svcCtx.DB.Task.UpdateOneID(res.ID).
					SetStatusStr("dispatch_failed").
					SetErrorMsg(fmt.Sprintf("Failed to dispatch task: %v", err)).
					Save(dispatchCtx)
				return
			}

			// 分配成功，更新状态为dispatched
			l.svcCtx.DB.Task.UpdateOneID(res.ID).
				SetStatusStr("dispatched").
				Save(dispatchCtx)

			l.Infow("Task dispatched successfully",
				logx.Field("task_id", res.TaskID),
				logx.Field("worker_id", selectedWorker.WorkerID))
		}()
	}

	return &newbee_ops_rpc.BaseIDResp{Id: res.ID, Msg: i18n.CreateSuccess}, nil
}
