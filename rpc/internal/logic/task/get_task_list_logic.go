package task

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/ent/task"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetTaskListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskListLogic {
	return &GetTaskListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskListLogic) GetTaskList(in *newbee_ops_rpc.TaskListReq) (*newbee_ops_rpc.TaskListResp, error) {
	var predicates []predicate.Task
	if in.CreatorId != nil {
		predicates = append(predicates, task.CreatorIDContains(in.GetCreatorId()))
	}
	if in.Executor != nil {
		predicates = append(predicates, task.ExecutorContains(in.GetExecutor()))
	}
	if in.StatusStr != nil {
		predicates = append(predicates, task.StatusStrContains(in.GetStatusStr()))
	}

	total, err := l.svcCtx.DB.Task.Query().Where(predicates...).Count(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}
	page, size := in.GetPage(), in.GetPageSize()
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	offset := int((page - 1) * size)
	list, err := l.svcCtx.DB.Task.Query().Where(predicates...).Limit(int(size)).Offset(offset).All(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &newbee_ops_rpc.TaskListResp{Total: uint64(total)}
	for _, v := range list {
		// 转换[]string到structpb.Struct (ci_ids)
		var ciIdsStruct *structpb.Struct
		if v.CiIds != nil {
			ciIdsAny := make(map[string]interface{})
			ciIdsList := make([]interface{}, len(v.CiIds))
			for i, val := range v.CiIds {
				ciIdsList[i] = val
			}
			ciIdsAny["list"] = ciIdsList
			var err error
			ciIdsStruct, err = structpb.NewStruct(ciIdsAny)
			if err != nil {
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
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
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
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
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
			}
		}

		resp.Data = append(resp.Data, &newbee_ops_rpc.TaskInfo{
			Id:             pointy.GetPointer(v.ID),
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
		})
	}
	return resp, nil
}
