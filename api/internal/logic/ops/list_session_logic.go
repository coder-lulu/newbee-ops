package ops

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionLogic {
	return &ListSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSessionLogic) ListSession(req *types.ListSessionReq) (resp *types.SessionListResp, err error) {
	// 1. 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100 // 最大100条
	}

	// 2. 构建查询条件
	listReq := &ops.SessionListReq{
		Page:     uint64(page),
		PageSize: uint64(size),
	}

	// 筛选条件
	if req.Status != "" {
		listReq.StatusStr = &req.Status
	}
	if req.CiId != "" {
		listReq.CiId = &req.CiId
	}
	if req.Protocol != "" {
		listReq.Protocol = &req.Protocol
	}
	if req.ProxyId != "" {
		listReq.ProxyId = &req.ProxyId
	}
	if req.UserId != "" {
		listReq.UserId = &req.UserId
	}

	// Note: Begin/End/SortBy/Order fields don't exist in SessionListReq protobuf
	// Sorting and time filtering is handled server-side automatically

	// 3. 查询session列表
	listResp, err := l.svcCtx.OpsClient.GetSessionList(l.ctx, listReq)
	if err != nil {
		l.Logger.Errorw("Failed to list sessions",
			logx.Field("page", page),
			logx.Field("size", size),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// 4. 转换为API响应格式
	items := make([]types.SessionItem, 0, len(listResp.Data))
	for _, session := range listResp.Data {
		items = append(items, types.SessionItem{
			Id:        session.GetSessionId(),
			TenantId:  "", // TenantId is managed by middleware, not exposed in protobuf
			UserId:    safeString(session.UserId),
			CiId:      safeString(session.CiId),
			Protocol:  safeString(session.Protocol),
			ProxyId:   safeString(session.ProxyId),
			Endpoint:  safeString(session.Endpoint),
			CreatedAt: safeInt64(session.CreatedAt),
			ExpiresAt: safeInt64(session.ExpiresAt),
			Status:    safeString(session.StatusStr),
			ClosedAt:  safeInt64(session.ClosedAt),
		})
	}

	l.Logger.Infow("Listed sessions",
		logx.Field("page", page),
		logx.Field("size", size),
		logx.Field("total", listResp.Total),
		logx.Field("count", len(items)))

	return &types.SessionListResp{
		Code: 0,
		Msg:  "success",
		Data: types.SessionListData{
			Items: items,
			Total: int(listResp.Total),
		},
	}, nil
}

// 辅助函数
