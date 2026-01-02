package task

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetTaskByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskByIdLogic {
	return &GetTaskByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskByIdLogic) GetTaskById(in *newbee_ops_rpc.IDReq) (*newbee_ops_rpc.TaskInfo, error) {
	result, err := l.svcCtx.DB.Task.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 转换[]string到structpb.Struct (ci_ids)
	var ciIdsStruct *structpb.Struct
	if result.CiIds != nil {
		ciIdsAny := make(map[string]interface{})
		list := make([]interface{}, len(result.CiIds))
		for i, v := range result.CiIds {
			list[i] = v
		}
		ciIdsAny["list"] = list
		var err error
		ciIdsStruct, err = structpb.NewStruct(ciIdsAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	// 转换map[string]string到structpb.Struct (command_env)
	var envStruct *structpb.Struct
	if result.CommandEnv != nil {
		envAny := make(map[string]interface{})
		for k, v := range result.CommandEnv {
			envAny[k] = v
		}
		var err error
		envStruct, err = structpb.NewStruct(envAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	// 转换map[string]string到structpb.Struct (tags)
	var tagsStruct *structpb.Struct
	if result.Tags != nil {
		tagsAny := make(map[string]interface{})
		for k, v := range result.Tags {
			tagsAny[k] = v
		}
		var err error
		tagsStruct, err = structpb.NewStruct(tagsAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	return &newbee_ops_rpc.TaskInfo{
		Id:             &result.ID,
		CreatedAt:      pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:      pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:         pointy.GetPointer(uint32(result.Status)),
		TaskId:         &result.TaskID,
		Name:           &result.Name,
		CreatorId:      &result.CreatorID,
		CiIds:          ciIdsStruct,
		Executor:       &result.Executor,
		CommandContent: &result.CommandContent,
		CommandTimeout: pointy.GetPointer(int64(result.CommandTimeout)),
		CommandWorkdir: &result.CommandWorkdir,
		CommandEnv:     envStruct,
		StatusStr:      &result.StatusStr,
		ResultOutput:   &result.ResultOutput,
		ErrorMsg:       &result.ErrorMsg,
		ExitCode:       pointy.GetPointer(int64(result.ExitCode)),
		RemoteTaskId:   &result.RemoteTaskID,
		StartTime:      &result.StartTime,
		EndTime:        &result.EndTime,
		Tags:           tagsStruct,
	}, nil
}
