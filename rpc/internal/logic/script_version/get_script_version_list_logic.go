package script_version

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/scriptversion"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetScriptVersionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScriptVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScriptVersionListLogic {
	return &GetScriptVersionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScriptVersionListLogic) GetScriptVersionList(in *ops.ScriptVersionListReq) (*ops.ScriptVersionListResp, error) {
	var predicates []predicate.ScriptVersion
	if in.CreatedAt != nil {
		predicates = append(predicates, scriptversion.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, scriptversion.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.ScriptId != nil {
		predicates = append(predicates, scriptversion.ScriptIDEQ(*in.ScriptId))
	}
	if in.Version != nil {
		predicates = append(predicates, scriptversion.VersionContains(*in.Version))
	}
	if in.Content != nil {
		predicates = append(predicates, scriptversion.ContentContains(*in.Content))
	}
	if in.Parameters != nil {
	}
	if in.ScriptType != nil {
		predicates = append(predicates, scriptversion.ScriptTypeContains(*in.ScriptType))
	}
	if in.Executor != nil {
		predicates = append(predicates, scriptversion.ExecutorContains(*in.Executor))
	}
	if in.ChangeLog != nil {
		predicates = append(predicates, scriptversion.ChangeLogContains(*in.ChangeLog))
	}
	if in.CreatedBy != nil {
		predicates = append(predicates, scriptversion.CreatedByEQ(*in.CreatedBy))
	}
	if in.Checksum != nil {
		predicates = append(predicates, scriptversion.ChecksumContains(*in.Checksum))
	}
	result, err := l.svcCtx.DB.ScriptVersion.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &ops.ScriptVersionListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &ops.ScriptVersionInfo{
			Id:          &v.ID,
			CreatedAt:   pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:   pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			ScriptId:	&v.ScriptID,
			Version:	&v.Version,
			Content:	&v.Content,
			Parameters: pointy.GetPointer(mustMarshalJSON(v.Parameters)),
			ScriptType:	&v.ScriptType,
			Executor:	&v.Executor,
			ChangeLog:	&v.ChangeLog,
			CreatedBy:	&v.CreatedBy,
			Checksum:	&v.Checksum,
		})
	}

	return resp, nil
}
