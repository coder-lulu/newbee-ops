package accessprofile

import (
	"context"

	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetAccessProfileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccessProfileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessProfileByIdLogic {
	return &GetAccessProfileByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccessProfileByIdLogic) GetAccessProfileById(in *newbee_ops_rpc.IDReq) (*newbee_ops_rpc.AccessProfileInfo, error) {
	result, err := l.svcCtx.DB.AccessProfile.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	// 转换[]string到structpb.Struct (capabilities)
	var capsStruct *structpb.Struct
	if result.Capabilities != nil {
		capsAny := make(map[string]interface{})
		list := make([]interface{}, len(result.Capabilities))
		for i, v := range result.Capabilities {
			list[i] = v
		}
		capsAny["list"] = list
		var err error
		capsStruct, err = structpb.NewStruct(capsAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	// 转换map[string]int到structpb.Struct (ports)
	var portsStruct *structpb.Struct
	if result.Ports != nil {
		portsAny := make(map[string]interface{})
		for k, v := range result.Ports {
			portsAny[k] = v
		}
		var err error
		portsStruct, err = structpb.NewStruct(portsAny)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	// 转换[]string到structpb.Struct (jump_chain)
	var jumpStruct *structpb.Struct
	if result.JumpChain != nil {
		jumpAny := make(map[string]interface{})
		list := make([]interface{}, len(result.JumpChain))
		for i, v := range result.JumpChain {
			list[i] = v
		}
		jumpAny["list"] = list
		var err error
		jumpStruct, err = structpb.NewStruct(jumpAny)
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

	return &newbee_ops_rpc.AccessProfileInfo{
		Id:            &result.ID,
		CreatedAt:     pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:     pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:        pointy.GetPointer(uint32(result.Status)),
		CiId:          &result.CiID,
		Capabilities:  capsStruct,
		Ports:         portsStruct,
		CredentialRef: &result.CredentialRef,
		PreferProxy:   &result.PreferProxy,
		JumpChain:     jumpStruct,
		Tags:          tagsStruct,
	}, nil
}
