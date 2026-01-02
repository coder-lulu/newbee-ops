package auditevent

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuditEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAuditEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuditEventLogic {
	return &CreateAuditEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAuditEventLogic) CreateAuditEvent(in *newbee_ops_rpc.AuditEventInfo) (*newbee_ops_rpc.BaseIDResp, error) {
	q := l.svcCtx.DB.AuditEvent.Create()
	if in.EventId != nil {
		q = q.SetEventID(in.GetEventId())
	}
	if in.Type != nil {
		q = q.SetType(in.GetType())
	}
	if in.SubjectUser != nil {
		q = q.SetSubjectUser(in.GetSubjectUser())
	}
	if in.SubjectCi != nil {
		q = q.SetSubjectCi(in.GetSubjectCi())
	}
	if in.SubjectProxy != nil {
		q = q.SetSubjectProxy(in.GetSubjectProxy())
	}
	if in.OccurAt != nil {
		q = q.SetOccurAt(in.GetOccurAt())
	}
	// 转换structpb.Struct到map[string]string
	if in.Meta != nil {
		metaAny := in.GetMeta().AsMap()
		meta := make(map[string]string)
		for k, v := range metaAny {
			if str, ok := v.(string); ok {
				meta[k] = str
			} else {
				l.Infow("skipping non-string meta value",
					logx.Field("key", k),
					logx.Field("type", fmt.Sprintf("%T", v)),
					logx.Field("value", v))
			}
		}
		if len(meta) > 0 {
			q = q.SetMeta(meta)
		}
	}
	if in.Status != nil {
		q = q.SetStatus(uint8(in.GetStatus()))
	}
	res, err := q.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}
	return &newbee_ops_rpc.BaseIDResp{Id: res.ID, Msg: i18n.CreateSuccess}, nil
}
