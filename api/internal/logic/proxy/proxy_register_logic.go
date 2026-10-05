package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy注册（PSK认证，无需JWT）
func NewProxyRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyRegisterLogic {
	return &ProxyRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyRegisterLogic) ProxyRegister(req *types.ProxyRegisterReq) (resp *types.BaseResp, err error) {
	ctx, rejection := registrationContext(l.ctx, l.svcCtx, req.PSK)
	if rejection != nil {
		return rejection, nil
	}
	l.ctx = ctx

	// 2. 检查是否已存在（通过 proxy_id）
	existing, _ := l.svcCtx.OpsClient.GetProxyByProxyId(l.ctx, req.ProxyID)

	// 3. 准备数据
	now := time.Now().UnixMilli()
	endpoints := req.Endpoints
	if endpoints == nil {
		endpoints = make(map[string]string)
	}

	capabilities := req.Capabilities
	if capabilities == nil {
		capabilities = []string{}
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	networkSegments := req.NetworkSegments
	if networkSegments == nil {
		networkSegments = []string{}
	}

	metadata := req.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	// 4. 创建或更新 Proxy
	if existing == nil {
		// 新注册
		proxyInfo := &ops.ProxyInfo{
			ProxyId:         &req.ProxyID,
			Name:            &req.Name,
			Ip:              &req.IP,
			Port:            pointy(int64(req.Port)),
			Version:         &req.Version,
			Region:          &req.Region,
			Zone:            &req.Zone,
			Capabilities:    pointy(mustMarshalJSON(capabilities)),
			Tags:            pointy(mustMarshalJSON(tags)),
			Endpoints:       pointy(mustMarshalJSON(endpoints)),
			ProxyStatus:     pointy("online"),
			RegisterTime:    pointy(now),
			LastHeartbeat:   pointy(now),
			MaxSessions:     pointy(int64(req.MaxSessions)),
			LocalIp:         &req.LocalIP,
			PublicIp:        &req.PublicIP,
			NetworkSegments: pointy(mustMarshalJSON(networkSegments)),
			Metadata:        pointy(mustMarshalJSON(metadata)),
			HealthCheckUrl:  &req.HealthCheckURL,
			Weight:          pointy(int64(100)), // 默认权重100
			Priority:        pointy(int64(0)),
			Status:          pointy(uint32(1)), // 启用
		}

		_, err = l.svcCtx.OpsClient.CreateProxy(l.ctx, proxyInfo)
		if err != nil {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("Failed to create proxy: %v", err),
			}, nil
		}

		l.Logger.Infow("Proxy registered successfully",
			logx.Field("proxy_id", req.ProxyID),
			logx.Field("name", req.Name),
			logx.Field("ip", req.IP))

	} else {
		// 更新已存在的记录
		proxyInfo := &ops.ProxyInfo{
			Id:              existing.Id,
			Name:            &req.Name,
			Ip:              &req.IP,
			Port:            pointy(int64(req.Port)),
			Version:         &req.Version,
			Region:          &req.Region,
			Zone:            &req.Zone,
			Capabilities:    pointy(mustMarshalJSON(capabilities)),
			Tags:            pointy(mustMarshalJSON(tags)),
			Endpoints:       pointy(mustMarshalJSON(endpoints)),
			ProxyStatus:     pointy("online"),
			LastHeartbeat:   pointy(now),
			MaxSessions:     pointy(int64(req.MaxSessions)),
			LocalIp:         &req.LocalIP,
			PublicIp:        &req.PublicIP,
			NetworkSegments: pointy(mustMarshalJSON(networkSegments)),
			Metadata:        pointy(mustMarshalJSON(metadata)),
			HealthCheckUrl:  &req.HealthCheckURL,
		}

		_, err = l.svcCtx.OpsClient.UpdateProxy(l.ctx, proxyInfo)
		if err != nil {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("Failed to update proxy: %v", err),
			}, nil
		}

		l.Logger.Infow("Proxy re-registered (updated)",
			logx.Field("proxy_id", req.ProxyID),
			logx.Field("name", req.Name))
	}

	return &types.BaseResp{
		Code: 0,
		Msg:  "success",
	}, nil
}

// pointy helper function
func pointy[T any](v T) *T { return &v }

// mustMarshalJSON marshals a value to JSON string, returns empty string on error
func mustMarshalJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
