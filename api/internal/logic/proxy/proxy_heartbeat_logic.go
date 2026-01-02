package proxy

import (
	"context"
	"fmt"
	"time"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyHeartbeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy心跳（PSK认证，无需JWT）
func NewProxyHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyHeartbeatLogic {
	return &ProxyHeartbeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyHeartbeatLogic) ProxyHeartbeat(req *types.ProxyHeartbeatReq) (resp *types.BaseResp, err error) {
	// 1. 验证 PSK
	if req.PSK != l.svcCtx.Config.Ops.Registration.PSK {
		return &types.BaseResp{
			Code: 401,
			Msg:  "Invalid PSK",
		}, nil
	}

	// 2. 查找 Proxy
	proxy, err := l.svcCtx.OpsClient.GetProxyByProxyId(l.ctx, req.ProxyID)
	if err != nil || proxy == nil {
		return &types.BaseResp{
			Code: 404,
			Msg:  "Proxy not found",
		}, nil
	}

	// 3. 更新心跳和指标
	now := time.Now().Unix()
	proxyInfo := &ops.ProxyInfo{
		Id:             proxy.Id,
		ProxyStatus:    &req.ProxyStatus,
		LastHeartbeat:  pointy(now),
		CpuUsage:       &req.CPUUsage,
		MemoryUsage:    &req.MemoryUsage,
		DiskUsage:      &req.DiskUsage,
		NetworkIn:      &req.NetworkIn,
		NetworkOut:     &req.NetworkOut,
		ActiveSessions: pointy(int64(req.ActiveSessions)),
		TotalRequests:  &req.TotalRequests,
		SuccessCount:   &req.SuccessCount,
		FailureCount:   &req.FailureCount,
	}

	_, err = l.svcCtx.OpsClient.UpdateProxy(l.ctx, proxyInfo)
	if err != nil {
		l.Logger.Errorw("Failed to update proxy heartbeat",
			logx.Field("proxy_id", req.ProxyID),
			logx.Field("error", err))
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to update proxy: %v", err),
		}, nil
	}

	// 4. 异步保存指标到历史表（可选）
	go l.saveMetrics(req, now)

	l.Logger.Debugw("Proxy heartbeat received",
		logx.Field("proxy_id", req.ProxyID),
		logx.Field("status", req.ProxyStatus),
		logx.Field("active_sessions", req.ActiveSessions))

	return &types.BaseResp{
		Code: 0,
		Msg:  "success",
	}, nil
}

// saveMetrics 保存指标到历史表（异步）
func (l *ProxyHeartbeatLogic) saveMetrics(req *types.ProxyHeartbeatReq, timestamp int64) {
	// 注意：这里使用新的 context，因为是异步操作
	ctx := context.Background()

	metricsInfo := &ops.ProxyMetricsInfo{
		ProxyId:         &req.ProxyID,
		Timestamp:       pointy(timestamp),
		CpuUsage:        &req.CPUUsage,
		MemoryUsage:     &req.MemoryUsage,
		DiskUsage:       &req.DiskUsage,
		NetworkInDelta:  &req.NetworkIn,
		NetworkOutDelta: &req.NetworkOut,
		ActiveSessions:  pointy(int64(req.ActiveSessions)),
		ProxyStatus:     &req.ProxyStatus,
	}

	// 这里需要调用 CreateProxyMetrics RPC 方法
	// 由于 OpsClient 接口还没有这个方法，我们先跳过
	// TODO: 添加 CreateProxyMetrics 到 OpsClient 接口

	_ = metricsInfo
	_ = ctx

	l.Logger.Debugw("Metrics saved (placeholder)",
		logx.Field("proxy_id", req.ProxyID))
}
