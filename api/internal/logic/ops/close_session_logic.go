package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloseSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseSessionLogic {
	return &CloseSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseSessionLogic) CloseSession(req *types.CloseSessionReq) (resp *types.CloseSessionResp, err error) {
	// 1. 查询session是否存在
	sessionInfo, err := l.svcCtx.OpsClient.GetSessionBySessionId(l.ctx, &ops.SessionSIDReq{
		SessionId: req.SessionId,
	})
	if err != nil {
		l.Logger.Errorw("Failed to get session",
			logx.Field("session_id", req.SessionId),
			logx.Field("error", err))
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 2. 检查session是否已经关闭
	if sessionInfo.StatusStr != nil && *sessionInfo.StatusStr == "closed" {
		l.Logger.Infow("Session already closed",
			logx.Field("session_id", req.SessionId),
			logx.Field("closed_at", *sessionInfo.ClosedAt))

		return &types.CloseSessionResp{
			Code: 0,
			Msg:  "session already closed",
			Data: types.CloseSessionData{
				Ok:       true,
				ClosedAt: *sessionInfo.ClosedAt,
			},
		}, nil
	}

	// 3. 更新session状态为closed
	now := time.Now().Unix()
	status := "closed"
	if req.Force {
		status = "force_closed"
	}

	updateReq := &ops.SessionInfo{
		Id:        sessionInfo.Id,
		StatusStr: &status,
		ClosedAt:  &now,
	}

	_, err = l.svcCtx.OpsClient.UpdateSession(l.ctx, updateReq)
	if err != nil {
		l.Logger.Errorw("Failed to update session status",
			logx.Field("session_id", req.SessionId),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to close session: %w", err)
	}

	l.Logger.Infow("Session closed successfully",
		logx.Field("session_id", req.SessionId),
		logx.Field("reason", req.Reason),
		logx.Field("force", req.Force),
		logx.Field("closed_at", now))

	return &types.CloseSessionResp{
		Code: 0,
		Msg:  "session closed successfully",
		Data: types.CloseSessionData{
			Ok:       true,
			ClosedAt: now,
		},
	}, nil
}
