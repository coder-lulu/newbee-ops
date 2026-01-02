package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/coder-lulu/newbee-ops-api/internal/client"
	proxylogic "github.com/coder-lulu/newbee-ops-api/internal/logic/proxy"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq) (resp *types.CreateSessionResp, err error) {
	// 1. 验证协议
	if !isValidProtocol(req.Protocol) {
		return nil, fmt.Errorf("invalid protocol: %s, must be one of: ssh, telnet, rdp, vnc", req.Protocol)
	}

	// 2. 选择合适的Proxy
	pickReq := &types.ProxyPickReq{
		Strategy:             "least_connections", // 默认使用最少连接数策略
		RequiredCapabilities: []string{req.Protocol},
	}

	// 如果指定了prefer_proxy，添加到过滤条件
	if req.PreferProxy != "" {
		pickReq.ExcludeProxyIDs = []string{} // 只使用指定的proxy
		// TODO: 这里需要根据PreferProxy筛选，暂时跳过
	}

	// 调用ProxyPick逻辑选择Proxy
	pickLogic := proxylogic.NewProxyPickLogic(l.ctx, l.svcCtx)
	proxyPickResp, err := pickLogic.ProxyPick(pickReq)
	if err != nil {
		l.Logger.Errorw("Failed to pick proxy",
			logx.Field("ci_id", req.CiId),
			logx.Field("protocol", req.Protocol),
			logx.Field("error", err))
		return nil, fmt.Errorf("no available proxy for protocol %s: %w", req.Protocol, err)
	}

	// 3. 获取连接凭证（从CI或AccessProfile）
	// TODO: 这里应该从CMDB或AccessProfile获取连接信息
	// 暂时使用Options中的参数作为连接信息
	target := req.Params["target"]
	username := req.Params["username"]
	password := req.Params["password"]

	if target == "" {
		return nil, fmt.Errorf("target host is required in params")
	}
	if username == "" {
		return nil, fmt.Errorf("username is required in params")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required in params")
	}

	// 默认端口
	port := req.Options.Port
	if port == 0 {
		port = getDefaultPort(req.Protocol)
	}

	// 4. 生成WebSocket URL
	proxyClient := client.NewProxyHTTPClient()
	var wsURL string

	// 构建HTTP endpoint（使用proxy的HTTP endpoint）
	proxyHTTPEndpoint := ""
	if httpEndpoint, ok := proxyPickResp.Data.Endpoints["http"]; ok {
		proxyHTTPEndpoint = httpEndpoint
	} else {
		// 如果没有http endpoint，使用IP:Port构建
		proxyHTTPEndpoint = fmt.Sprintf("http://%s:%d", proxyPickResp.Data.IP, proxyPickResp.Data.Port)
	}

	switch req.Protocol {
	case "ssh":
		wsURL = proxyClient.GenerateSSHWebSocketURL(proxyHTTPEndpoint, target, port, username, password)
	case "telnet":
		wsURL = proxyClient.GenerateTelnetWebSocketURL(proxyHTTPEndpoint, target, port, username, password)
	case "rdp":
		wsURL = proxyClient.GenerateRDPWebSocketURL(proxyHTTPEndpoint, target, port, username, password)
	case "vnc":
		wsURL = proxyClient.GenerateVNCWebSocketURL(proxyHTTPEndpoint, target, port, password)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", req.Protocol)
	}

	// 5. 创建Session记录到数据库
	now := time.Now().Unix()
	expiresAt := now + 7200 // 默认2小时过期

	sessionInfo := &ops.SessionInfo{
		CiId:      &req.CiId,
		Protocol:  &req.Protocol,
		ProxyId:   &proxyPickResp.Data.ProxyID,
		Endpoint:  &proxyHTTPEndpoint,
		CreatedAt: pointy(now),
		ExpiresAt: pointy(expiresAt),
		StatusStr: pointy("active"),
	}

	createResp, err := l.svcCtx.OpsClient.CreateSession(l.ctx, sessionInfo)
	if err != nil {
		l.Logger.Errorw("Failed to create session record",
			logx.Field("ci_id", req.CiId),
			logx.Field("protocol", req.Protocol),
			logx.Field("proxy_id", proxyPickResp.Data.ProxyID),
			logx.Field("error", err))
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 6. 生成Session ID（使用数据库返回的ID）
	sessionID := fmt.Sprintf("session_%d", createResp.Id)

	// 7. 生成访问Token（可选，用于WebSocket鉴权）
	// TODO: 这里应该生成一个临时token，用于WebSocket连接鉴权
	token := fmt.Sprintf("token_%d_%d", createResp.Id, now)

	// 8. 记录握手参数
	handshake := make(map[string]string)
	handshake["protocol"] = req.Protocol
	handshake["proxy_id"] = proxyPickResp.Data.ProxyID
	handshake["target"] = target
	if req.Options.Width > 0 {
		handshake["width"] = fmt.Sprintf("%d", req.Options.Width)
	}
	if req.Options.Height > 0 {
		handshake["height"] = fmt.Sprintf("%d", req.Options.Height)
	}

	l.Logger.Infow("Session created successfully",
		logx.Field("session_id", sessionID),
		logx.Field("ci_id", req.CiId),
		logx.Field("protocol", req.Protocol),
		logx.Field("proxy_id", proxyPickResp.Data.ProxyID),
		logx.Field("ws_url", wsURL))

	return &types.CreateSessionResp{
		Code: 0,
		Msg:  "session created successfully",
		Data: types.CreateSessionData{
			SessionId: sessionID,
			ProxyId:   proxyPickResp.Data.ProxyID,
			WsUrl:     wsURL,
			Token:     token,
			ExpiresAt: expiresAt,
			Handshake: handshake,
		},
	}, nil
}

// isValidProtocol 验证协议是否有效
func isValidProtocol(protocol string) bool {
	validProtocols := []string{"ssh", "telnet", "rdp", "vnc"}
	for _, p := range validProtocols {
		if p == protocol {
			return true
		}
	}
	return false
}

// getDefaultPort 获取协议默认端口
func getDefaultPort(protocol string) int {
	switch protocol {
	case "ssh":
		return 22
	case "telnet":
		return 23
	case "rdp":
		return 3389
	case "vnc":
		return 5900
	default:
		return 22
	}
}

// pointy helper function
func pointy[T any](v T) *T { return &v }

// safeString returns empty string if pointer is nil
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// safeInt64 returns 0 if pointer is nil
func safeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// safeUint64ToString converts uint64 pointer to string
func safeUint64ToString(i *uint64) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}
