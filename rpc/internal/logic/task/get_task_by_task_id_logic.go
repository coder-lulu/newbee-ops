package task

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/ent/task"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetTaskByTaskIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskByTaskIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskByTaskIdLogic {
	return &GetTaskByTaskIdLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetTaskByTaskIdLogic) GetTaskByTaskId(in *newbee_ops_rpc.TaskIdReq) (*newbee_ops_rpc.TaskInfo, error) {
	v, err := l.svcCtx.DB.Task.Query().Where(task.TaskIDEQ(in.GetTaskId())).Only(l.ctx)
	if err != nil {
		return nil, err
	}

	// 转换[]string到structpb.Struct (ci_ids)
	var ciIdsStruct *structpb.Struct
	if v.CiIds != nil {
		ciIdsAny := make(map[string]interface{})
		list := make([]interface{}, len(v.CiIds))
		for i, val := range v.CiIds {
			list[i] = val
		}
		ciIdsAny["list"] = list
		var err error
		ciIdsStruct, err = structpb.NewStruct(ciIdsAny)
		if err != nil {
			return nil, err
		}
	}

	// 转换map[string]string到structpb.Struct (command_env)
	var envStruct *structpb.Struct
	if v.CommandEnv != nil {
		envAny := make(map[string]interface{})
		for k, val := range v.CommandEnv {
			envAny[k] = val
		}
		var err error
		envStruct, err = structpb.NewStruct(envAny)
		if err != nil {
			return nil, err
		}
	}

	// 转换map[string]string到structpb.Struct (tags)
	var tagsStruct *structpb.Struct
	if v.Tags != nil {
		tagsAny := make(map[string]interface{})
		for k, val := range v.Tags {
			tagsAny[k] = val
		}
		var err error
		tagsStruct, err = structpb.NewStruct(tagsAny)
		if err != nil {
			return nil, err
		}
	}

	return &newbee_ops_rpc.TaskInfo{
		Id:             &v.ID,
		CreatedAt:      pointy.GetPointer(v.CreatedAt.UnixMilli()),
		UpdatedAt:      pointy.GetPointer(v.UpdatedAt.UnixMilli()),
		Status:         pointy.GetPointer(uint32(v.Status)),
		TaskId:         &v.TaskID,
		Name:           &v.Name,
		CreatorId:      &v.CreatorID,
		CiIds:          ciIdsStruct,
		Executor:       &v.Executor,
		CommandContent: &v.CommandContent,
		CommandTimeout: pointy.GetPointer(int64(v.CommandTimeout)),
		CommandWorkdir: &v.CommandWorkdir,
		CommandEnv:     envStruct,
		StatusStr:      &v.StatusStr,
		ResultOutput:   &v.ResultOutput,
		ErrorMsg:       &v.ErrorMsg,
		ExitCode:       pointy.GetPointer(int64(v.ExitCode)),
		RemoteTaskId:   &v.RemoteTaskID,
		StartTime:      &v.StartTime,
		EndTime:        &v.EndTime,
		Tags:           tagsStruct,
	}, nil
}
