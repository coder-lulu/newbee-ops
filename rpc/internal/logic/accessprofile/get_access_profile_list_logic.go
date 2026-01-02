package accessprofile

import (
	"context"
	"time"

	"github.com/coder-lulu/newbee-ops-rpc/ent/accessprofile"
	"github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
	"github.com/coder-lulu/newbee-ops-rpc/internal/svc"
	"github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
	newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-common/v2/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/structpb"
)

type GetAccessProfileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccessProfileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessProfileListLogic {
	return &GetAccessProfileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccessProfileListLogic) GetAccessProfileList(in *newbee_ops_rpc.AccessProfileListReq) (*newbee_ops_rpc.AccessProfileListResp, error) {
	var predicates []predicate.AccessProfile
	if in.CreatedAt != nil {
		predicates = append(predicates, accessprofile.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, accessprofile.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, accessprofile.StatusEQ(uint8(*in.Status)))
	}
	if in.CiId != nil {
		predicates = append(predicates, accessprofile.CiIDContains(*in.CiId))
	}
    // skip JSON filters in generic impl
	if in.CredentialRef != nil {
		predicates = append(predicates, accessprofile.CredentialRefContains(*in.CredentialRef))
	}
	if in.PreferProxy != nil {
		predicates = append(predicates, accessprofile.PreferProxyContains(*in.PreferProxy))
	}
	// skip JSON filters in generic impl
	total, err := l.svcCtx.DB.AccessProfile.Query().Where(predicates...).Count(l.ctx)
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
	list, err := l.svcCtx.DB.AccessProfile.Query().Where(predicates...).Limit(int(size)).Offset(offset).All(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &newbee_ops_rpc.AccessProfileListResp{Total: uint64(total)}
	for _, v := range list {
		// 转换[]string到structpb.Struct (capabilities)
		var capsStruct *structpb.Struct
		if v.Capabilities != nil {
			capsAny := make(map[string]interface{})
			capList := make([]interface{}, len(v.Capabilities))
			for i, val := range v.Capabilities {
				capList[i] = val
			}
			capsAny["list"] = capList
			var err error
			capsStruct, err = structpb.NewStruct(capsAny)
			if err != nil {
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
			}
		}

		// 转换map[string]int到structpb.Struct (ports)
		var portsStruct *structpb.Struct
		if v.Ports != nil {
			portsAny := make(map[string]interface{})
			for k, val := range v.Ports {
				portsAny[k] = val
			}
			var err error
			portsStruct, err = structpb.NewStruct(portsAny)
			if err != nil {
				return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
			}
		}

		// 转换[]string到structpb.Struct (jump_chain)
		var jumpStruct *structpb.Struct
		if v.JumpChain != nil {
			jumpAny := make(map[string]interface{})
			jumpList := make([]interface{}, len(v.JumpChain))
			for i, val := range v.JumpChain {
				jumpList[i] = val
			}
			jumpAny["list"] = jumpList
			var err error
			jumpStruct, err = structpb.NewStruct(jumpAny)
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

		resp.Data = append(resp.Data, &newbee_ops_rpc.AccessProfileInfo{
			Id:            &v.ID,
			CreatedAt:     pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:     pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:        pointy.GetPointer(uint32(v.Status)),
			CiId:          &v.CiID,
			Capabilities:  capsStruct,
			Ports:         portsStruct,
			CredentialRef: &v.CredentialRef,
			PreferProxy:   &v.PreferProxy,
			JumpChain:     jumpStruct,
			Tags:          tagsStruct,
		})
	}
	return resp, nil
}
