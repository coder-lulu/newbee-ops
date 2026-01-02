package credentialref

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCredentialRefLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCredentialRefLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCredentialRefLogic {
	return &CreateCredentialRefLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCredentialRefLogic) CreateCredentialRef(in *newbee_ops_rpc.CredentialRefInfo) (*newbee_ops_rpc.BaseIDResp, error) {
	q := l.svcCtx.DB.CredentialRef.Create()
	if in.Provider != nil {
		q = q.SetProvider(in.GetProvider())
	}
	if in.Ref != nil {
		q = q.SetRef(in.GetRef())
	}
	if in.Scope != nil {
		q = q.SetScope(in.GetScope())
	}
	if in.CreatedBy != nil {
		q = q.SetCreatedBy(in.GetCreatedBy())
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
