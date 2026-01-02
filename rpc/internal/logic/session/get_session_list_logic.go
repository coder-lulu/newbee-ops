package session

import (
    "context"
    "time"

    "github.com/coder-lulu/newbee-ops-rpc/ent/predicate"
    "github.com/coder-lulu/newbee-ops-rpc/ent/session"
    "github.com/coder-lulu/newbee-ops-rpc/internal/svc"
    "github.com/coder-lulu/newbee-ops-rpc/internal/utils/dberrorhandler"
    newbee_ops_rpc "github.com/coder-lulu/newbee-ops-rpc/types/ops"

    "github.com/coder-lulu/newbee-common/v2/utils/pointy"
    "github.com/zeromicro/go-zero/core/logx"
    "google.golang.org/protobuf/types/known/structpb"
)

type GetSessionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionListLogic) GetSessionList(in *newbee_ops_rpc.SessionListReq) (*newbee_ops_rpc.SessionListResp, error) {
	var predicates []predicate.Session
	if in.CreatedAt != nil {
		predicates = append(predicates, session.CreatedAtGTE(time.UnixMilli(*in.CreatedAt)))
	}
	if in.UpdatedAt != nil {
		predicates = append(predicates, session.UpdatedAtGTE(time.UnixMilli(*in.UpdatedAt)))
	}
	if in.Status != nil {
		predicates = append(predicates, session.StatusEQ(uint8(*in.Status)))
	}
	if in.SessionId != nil {
		predicates = append(predicates, session.SessionIDContains(*in.SessionId))
	}
	if in.UserId != nil {
		predicates = append(predicates, session.UserIDContains(*in.UserId))
	}
	if in.CiId != nil {
		predicates = append(predicates, session.CiIDContains(*in.CiId))
	}
	if in.Protocol != nil {
		predicates = append(predicates, session.ProtocolContains(*in.Protocol))
	}
	if in.ProxyId != nil {
		predicates = append(predicates, session.ProxyIDContains(*in.ProxyId))
	}
	if in.StatusStr != nil {
		predicates = append(predicates, session.StatusStrContains(*in.StatusStr))
	}
	// skip JSON tags filter in generic impl

	total, err := l.svcCtx.DB.Session.Query().Where(predicates...).Count(l.ctx)
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
	list, err := l.svcCtx.DB.Session.Query().Where(predicates...).Limit(int(size)).Offset(offset).All(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &newbee_ops_rpc.SessionListResp{Total: uint64(total)}
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

		resp.Data = append(resp.Data, &newbee_ops_rpc.SessionInfo{
			Id:        &v.ID,
			CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt: pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Status:    pointy.GetPointer(uint32(v.Status)),
			SessionId: &v.SessionID,
			UserId:    &v.UserID,
			CiId:      &v.CiID,
			Protocol:  &v.Protocol,
			ProxyId:   &v.ProxyID,
			Endpoint:  &v.Endpoint,
			StatusStr: &v.StatusStr,
			ExpiresAt: &v.ExpiresAt,
			ClosedAt:  &v.ClosedAt,
			Tags:      tagsStruct,
		})
	}
	return resp, nil
}
