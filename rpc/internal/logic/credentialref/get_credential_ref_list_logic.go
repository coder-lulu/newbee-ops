package credentialref

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/credentialref"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetCredentialRefListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCredentialRefListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCredentialRefListLogic {
	return &GetCredentialRefListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCredentialRefListLogic) GetCredentialRefList(in *newbee_ops_rpc.CredentialRefListReq) (*newbee_ops_rpc.CredentialRefListResp, error) {
	var predicates []predicate.CredentialRef
	if in.CreatedAt != nil {
		predicates = append(predicates, credentialref.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, credentialref.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, credentialref.StatusEQ(uint8(*in.Status)))
	}
	if in.Provider != nil {
		predicates = append(predicates, credentialref.ProviderContains(*in.Provider))
	}
	if in.Ref != nil {
		predicates = append(predicates, credentialref.RefContains(*in.Ref))
	}
	if in.Scope != nil {
		predicates = append(predicates, credentialref.ScopeContains(*in.Scope))
	}
	if in.CreatedBy != nil {
		predicates = append(predicates, credentialref.CreatedByContains(*in.CreatedBy))
	}
    // skip JSON tags filter in generic impl
    total, err := l.svcCtx.DB.CredentialRef.Query().Where(predicates...).Count(l.ctx)
    if err != nil { return nil, dberrorhandler.DefaultEntError(l.Logger, err, in) }
    page, size := in.GetPage(), in.GetPageSize(); if page == 0 { page = 1 }; if size == 0 { size = 10 }
    offset := int((page-1)*size)
    list, err := l.svcCtx.DB.CredentialRef.Query().Where(predicates...).Limit(int(size)).Offset(offset).All(l.ctx)
    if err != nil { return nil, dberrorhandler.DefaultEntError(l.Logger, err, in) }

	resp := &newbee_ops_rpc.CredentialRefListResp{Total: uint64(total)}
	for _, v := range list {
		// 转换map[string]string到structpb.Struct
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

		resp.Data = append(resp.Data, &newbee_ops_rpc.CredentialRefInfo{
			Id:        &v.ID,
			CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt: pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:    pointy.GetPointer(uint32(v.Status)),
			Provider:  &v.Provider,
			Ref:       &v.Ref,
			Scope:     &v.Scope,
			CreatedBy: &v.CreatedBy,
			Tags:      tagsStruct,
		})
	}
	return resp, nil
}
