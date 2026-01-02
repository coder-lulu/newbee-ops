package auditevent

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/auditevent"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetAuditEventListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAuditEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditEventListLogic {
	return &GetAuditEventListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAuditEventListLogic) GetAuditEventList(in *newbee_ops_rpc.AuditEventListReq) (*newbee_ops_rpc.AuditEventListResp, error) {
	var predicates []predicate.AuditEvent
	if in.CreatedAt != nil {
		predicates = append(predicates, auditevent.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, auditevent.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, auditevent.StatusEQ(uint8(*in.Status)))
	}
	if in.EventId != nil {
		predicates = append(predicates, auditevent.EventIDContains(*in.EventId))
	}
	if in.Type != nil {
		predicates = append(predicates, auditevent.TypeContains(*in.Type))
	}
	if in.SubjectUser != nil {
		predicates = append(predicates, auditevent.SubjectUserContains(*in.SubjectUser))
	}
	if in.SubjectCi != nil {
		predicates = append(predicates, auditevent.SubjectCiContains(*in.SubjectCi))
	}
	if in.SubjectProxy != nil {
		predicates = append(predicates, auditevent.SubjectProxyContains(*in.SubjectProxy))
	}
	if in.OccurAt != nil {
		predicates = append(predicates, auditevent.OccurAtEQ(*in.OccurAt))
	}
    total, err := l.svcCtx.DB.AuditEvent.Query().Where(predicates...).Count(l.ctx)
    if err != nil { return nil, dberrorhandler.DefaultEntError(l.Logger, err, in) }
    page, size := in.GetPage(), in.GetPageSize(); if page == 0 { page = 1 }; if size == 0 { size = 10 }
    offset := int((page-1)*size)
    list, err := l.svcCtx.DB.AuditEvent.Query().Where(predicates...).Limit(int(size)).Offset(offset).All(l.ctx)
    if err != nil { return nil, dberrorhandler.DefaultEntError(l.Logger, err, in) }

	resp := &newbee_ops_rpc.AuditEventListResp{Total: uint64(total)}
	for _, v := range list {
		// 转换map[string]string到structpb.Struct
		var metaStruct *structpb.Struct
		if v.Meta != nil {
			metaAny := make(map[string]interface{})
			for k, val := range v.Meta {
				metaAny[k] = val
			}
			var err error
			metaStruct, err = structpb.NewStruct(metaAny)
			if err != nil {
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
			}
		}

		resp.Data = append(resp.Data, &newbee_ops_rpc.AuditEventInfo{
			Id:           &v.ID,
			CreatedAt:    pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:    pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:       pointy.GetPointer(uint32(v.Status)),
			EventId:      &v.EventID,
			Type:         &v.Type,
			SubjectUser:  &v.SubjectUser,
			SubjectCi:    &v.SubjectCi,
			SubjectProxy: &v.SubjectProxy,
			OccurAt:      &v.OccurAt,
			Meta:         metaStruct,
		})
	}
	return resp, nil
}
