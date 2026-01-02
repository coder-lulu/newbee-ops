package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/coder-lulu/newbee-ops-api/internal/client"
	proxylogic "github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.CreateTaskResp, err error) {
	// 1. 验证任务参数
	if len(req.CiIds) == 0 {
		return nil, fmt.Errorf("ci_ids is required")
	}
	if req.Command.Content == "" {
		return nil, fmt.Errorf("command content is required")
	}
	if req.Executor == "" {
		req.Executor = "ssh" // 默认使用ssh执行器
	}

	// 设置默认超时时间
	if req.Command.Timeout == 0 {
		req.Command.Timeout = 300 // 默认5分钟超时
	}

	// 2. 生成任务ID
	now := time.Now().Unix()
	taskID := fmt.Sprintf("task_%d", now)

	// 3. 创建Task记录到数据库
	taskInfo := &ops.TaskInfo{
		TaskId:    &taskID,
		Name:      &req.Name,
		StatusStr: pointy("pending"),
		CreatedAt: pointy(now),
	}

	createResp, err := l.svcCtx.OpsClient.CreateTask(l.ctx, taskInfo)
	if err != nil {
		l.Logger.Errorw("Failed to create task record",
			logx.Field("task_id", taskID),
			logx.Field("name", req.Name),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	l.Logger.Infow("Task record created",
		logx.Field("task_id", taskID),
		logx.Field("db_id", createResp.Id),
		logx.Field("name", req.Name),
		logx.Field("ci_count", len(req.CiIds)))

	// 4. 为每个CI下发任务
	proxyClient := client.NewProxyHTTPClient()
	successCount := 0
	failedCIs := make([]string, 0)

	for _, ciId := range req.CiIds {
		// 4.1 选择Proxy
		pickReq := &types.ProxyPickReq{
			Strategy:             "least_connections",
			RequiredCapabilities: []string{req.Executor}, // ssh or agent
		}

		pickLogic := proxylogic.NewProxyPickLogic(l.ctx, l.svcCtx)
		proxyPickResp, err := pickLogic.ProxyPick(pickReq)
		if err != nil {
			l.Logger.Errorw("Failed to pick proxy for CI",
				logx.Field("ci_id", ciId),
				logx.Field("error", err))
			failedCIs = append(failedCIs, ciId)
			continue
		}

		// 4.2 构建Proxy HTTP endpoint
		proxyHTTPEndpoint := ""
		if httpEndpoint, ok := proxyPickResp.Data.Endpoints["http"]; ok {
			proxyHTTPEndpoint = httpEndpoint
		} else {
			proxyHTTPEndpoint = fmt.Sprintf("http://%s:%d", proxyPickResp.Data.IP, proxyPickResp.Data.Port)
		}

		// 4.3 根据执行器类型下发任务
		var taskResp *client.TaskResponse

		if req.Executor == "ssh" || req.Executor == "telnet" {
			// SSH/Telnet执行器：执行命令
			// TODO: 这里应该从CMDB获取CI的连接信息
			// 暂时跳过实际的任务下发，因为需要CI的连接信息

			cmdReq := &client.CommandExecuteRequest{
				Target:      "placeholder", // TODO: 从CMDB获取
				Port:        22,
				Protocol:    req.Executor,
				Username:    "placeholder", // TODO: 从CMDB获取
				Password:    "placeholder", // TODO: 从CMDB获取
				Command:     req.Command.Content,
				WorkingDir:  req.Command.WorkDir,
				Timeout:     int32(req.Command.Timeout),
				Environment: req.Command.Env,
			}

			taskResp, err = proxyClient.ExecuteCommand(l.ctx, proxyHTTPEndpoint, cmdReq)
		} else if req.Executor == "agent" {
			// Agent执行器：通过脚本执行
			scriptReq := &client.ScriptExecuteRequest{
				Target:        "placeholder", // TODO: 从CMDB获取
				Port:          22,
				Protocol:      "ssh",
				Username:      "placeholder",
				Password:      "placeholder",
				ScriptContent: req.Command.Content,
				ScriptType:    "bash",
				WorkingDir:    req.Command.WorkDir,
				Timeout:       int32(req.Command.Timeout),
				Environment:   req.Command.Env,
				CleanupAfter:  true,
			}

			taskResp, err = proxyClient.ExecuteScript(l.ctx, proxyHTTPEndpoint, scriptReq)
		} else {
			l.Logger.Errorw("Unsupported executor type",
				logx.Field("executor", req.Executor),
				logx.Field("ci_id", ciId))
			failedCIs = append(failedCIs, ciId)
			continue
		}

		if err != nil {
			l.Logger.Errorw("Failed to dispatch task to proxy",
				logx.Field("ci_id", ciId),
				logx.Field("proxy_id", proxyPickResp.Data.ProxyID),
				logx.Field("error", err))
			failedCIs = append(failedCIs, ciId)
			continue
		}

		l.Logger.Infow("Task dispatched to proxy",
			logx.Field("task_id", taskID),
			logx.Field("proxy_task_id", taskResp.TaskID),
			logx.Field("ci_id", ciId),
			logx.Field("proxy_id", proxyPickResp.Data.ProxyID))

		successCount++
	}

	// 5. 更新任务状态
	status := "running"
	if successCount == 0 {
		status = "failed"
	} else if len(failedCIs) > 0 {
		status = "partial"
	}

	// TODO: 更新任务状态到数据库
	l.Logger.Infow("Task dispatch completed",
		logx.Field("task_id", taskID),
		logx.Field("total_cis", len(req.CiIds)),
		logx.Field("success_count", successCount),
		logx.Field("failed_count", len(failedCIs)),
		logx.Field("status", status))

	return &types.CreateTaskResp{
		Code: 0,
		Msg:  "task created successfully",
		Data: types.CreateTaskData{
			TaskId: taskID,
			Status: status,
		},
	}, nil
}

// pointy helper function
