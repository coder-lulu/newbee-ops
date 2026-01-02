package ops

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionByQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSessionByQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionByQueryLogic {
	return &GetSessionByQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionByQueryLogic) GetSessionByQuery(req *types.GetSessionQueryReq) (resp *types.SessionItem, err error) {
	// 1. 如果提供了ID，直接查询
	if req.Id != "" {
		return NewGetSessionLogic(l.ctx, l.svcCtx).GetSession(&types.GetSessionReq{
			Id: req.Id,
		})
	}

	// 2. 构建查询条件
	listReq := &ops.SessionListReq{
		Page:     1,
		PageSize: 1, // 只需要第一条
	}

	if req.CiId != "" {
		listReq.CiId = &req.CiId
	}
	if req.Status != "" {
		listReq.StatusStr = &req.Status
	}
	if req.Protocol != "" {
		listReq.Protocol = &req.Protocol
	}
	if req.ProxyId != "" {
		listReq.ProxyId = &req.ProxyId
	}
	if req.UserId != "" {
		// 注意：UserId是字符串，需要转换
		listReq.UserId = &req.UserId
	}

	// 3. 查询session列表
	listResp, err := l.svcCtx.OpsClient.GetSessionList(l.ctx, listReq)
	if err != nil {
		l.Logger.Errorw("Failed to query sessions",
			logx.Field("ci_id", req.CiId),
			logx.Field("status", req.Status),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}

	// 4. 检查是否有结果
	if listResp.Total == 0 || len(listResp.Data) == 0 {
		return nil, fmt.Errorf("no session found matching criteria")
	}

	// 5. 返回第一个匹配的session
	session := listResp.Data[0]

	return &types.SessionItem{
		Id:        fmt.Sprintf("session_%d", *session.Id),
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
	}, nil
}

// 辅助函数（复用GetSessionLogic中的函数）


