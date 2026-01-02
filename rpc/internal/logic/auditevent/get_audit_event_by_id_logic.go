package auditevent

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetAuditEventByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuditEventByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditEventByIdLogic {
	return &GetAuditEventByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuditEventByIdLogic) GetAuditEventById(in *newbee_ops_rpc.IDReq) (*newbee_ops_rpc.AuditEventInfo, error) {
	result, err := l.svcCtx.DB.AuditEvent.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 转换map[string]string到structpb.Struct
	var metaStruct *structpb.Struct
	if result.Meta != nil {
		metaAny := make(map[string]interface{})
		for k, v := range result.Meta {
			metaAny[k] = v
		}
		var err error
		metaStruct, err = structpb.NewStruct(metaAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	return &newbee_ops_rpc.AuditEventInfo{
		Id:           &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:       pointy.GetPointer(uint32(result.Status)),
		EventId:      &result.EventID,
		Type:         &result.Type,
		SubjectUser:  &result.SubjectUser,
		SubjectCi:    &result.SubjectCi,
		SubjectProxy: &result.SubjectProxy,
		OccurAt:      &result.OccurAt,
		Meta:         metaStruct,
	}, nil
}
