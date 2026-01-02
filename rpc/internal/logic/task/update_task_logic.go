package task

import (
	"context"
	"fmt"
	"math"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTaskLogic) UpdateTask(in *newbee_ops_rpc.TaskInfo) (*newbee_ops_rpc.BaseResp, error) {
	q := l.svcCtx.DB.Task.UpdateOneID(in.GetId())
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
	if err := q.Exec(l.ctx); err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}
	return &newbee_ops_rpc.BaseResp{Msg: i18n.UpdateSuccess}, nil
}
