package ops

import (
	"context"
	"fmt"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
	return &GetSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionLogic) GetSession(req *types.GetSessionReq) (resp *types.SessionItem, err error) {
	// 1. 从RPC查询session (session_id格式: "session_12345" 或直接数字ID)
	sessionInfo, err := l.svcCtx.OpsClient.GetSessionBySessionId(l.ctx, &ops.SessionSIDReq{
		SessionId: req.Id,
	})
	if err != nil {
		l.Logger.Errorw("Failed to get session",
			logx.Field("session_id", req.Id),
			logx.Field("error", err))
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 2. 转换为SessionItem
	return &types.SessionItem{
		Id:        fmt.Sprintf("session_%d", *sessionInfo.Id),
		TenantId:  "", // TenantId is managed by middleware, not exposed in protobuf
		UserId:    safeString(sessionInfo.UserId),
		CiId:      safeString(sessionInfo.CiId),
		Protocol:  safeString(sessionInfo.Protocol),
		ProxyId:   safeString(sessionInfo.ProxyId),
		Endpoint:  safeString(sessionInfo.Endpoint),
		CreatedAt: safeInt64(sessionInfo.CreatedAt),
		ExpiresAt: safeInt64(sessionInfo.ExpiresAt),
		Status:    safeString(sessionInfo.StatusStr),
		ClosedAt:  safeInt64(sessionInfo.ClosedAt),
	}, nil
}

// 辅助函数：安全地转换指针


