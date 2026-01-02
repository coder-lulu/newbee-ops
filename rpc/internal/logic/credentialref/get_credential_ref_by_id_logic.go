package credentialref

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetCredentialRefByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCredentialRefByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCredentialRefByIdLogic {
	return &GetCredentialRefByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCredentialRefByIdLogic) GetCredentialRefById(in *newbee_ops_rpc.IDReq) (*newbee_ops_rpc.CredentialRefInfo, error) {
	result, err := l.svcCtx.DB.CredentialRef.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 转换map[string]string到structpb.Struct
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

	return &newbee_ops_rpc.CredentialRefInfo{
		Id:        &result.ID,
		CreatedAt: pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt: pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:    pointy.GetPointer(uint32(result.Status)),
		Provider:  &result.Provider,
		Ref:       &result.Ref,
		Scope:     &result.Scope,
		CreatedBy: &result.CreatedBy,
		Tags:      tagsStruct,
	}, nil
}
