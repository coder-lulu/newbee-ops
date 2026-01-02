package session

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *newbee_ops_rpc.SessionInfo) (*newbee_ops_rpc.BaseIDResp, error) {
	q := l.svcCtx.DB.Session.Create()
	if in.SessionId != nil {
		q = q.SetSessionID(in.GetSessionId())
	}
	if in.UserId != nil {
		q = q.SetUserID(in.GetUserId())
	}
	if in.CiId != nil {
		q = q.SetCiID(in.GetCiId())
	}
	if in.Protocol != nil {
		q = q.SetProtocol(in.GetProtocol())
	}
	if in.ProxyId != nil {
		q = q.SetProxyID(in.GetProxyId())
	}
	if in.Endpoint != nil {
		q = q.SetEndpoint(in.GetEndpoint())
	}
	if in.StatusStr != nil {
		q = q.SetStatusStr(in.GetStatusStr())
	}
	if in.ExpiresAt != nil {
		q = q.SetExpiresAt(in.GetExpiresAt())
	}
	if in.ClosedAt != nil {
		q = q.SetClosedAt(in.GetClosedAt())
	}
	// 转换structpb.Struct到map[string]string
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
	res, err := q.Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}
	return &newbee_ops_rpc.BaseIDResp{Id: res.ID, Msg: i18n.CreateSuccess}, nil
}
