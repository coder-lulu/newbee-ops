package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"

	"github.com/coder-lulu/newbee-ops-api/internal/domain"
	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"github.com/golang-jwt/jwt/v4"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateSessionLogic) Create(req *types.CreateSessionReq) (*types.CreateSessionResp, error) {
	// load profile and validate protocol (if present)
	prof, ok, err := domain.LoadProfile(l.ctx, l.svcCtx, req.CiId)
	if err != nil {
		return nil, err
	}
	if ok && !domain.AllowProtocol(prof, req.Protocol) {
		return nil, errBadRequest("protocol not allowed by profile")
	}
	// Phase 3: 使用WorkerClient选择Worker（替代旧的Proxy选择）
	var endpoint, proxyId, workerID, workerIP string
	var workerPort int

	if l.svcCtx.WorkerClient != nil {
		// 使用Worker选择逻辑
		worker, err := l.svcCtx.WorkerClient.PickWorker(l.ctx, []string{req.Protocol}, "", nil)
		if err != nil {
			return nil, errBadRequest("no available worker for protocol: " + req.Protocol)
		}

		workerID = worker.WorkerID
		workerIP = worker.IP
		workerPort = worker.Port

		// 从Worker的endpoints构建WebSocket endpoint
		if wsEndpoint, ok := worker.Endpoints["ws"]; ok {
			endpoint = wsEndpoint
		} else if httpEndpoint, ok := worker.Endpoints["http"]; ok {
			// 如果没有ws endpoint，尝试从http构建
			endpoint = httpEndpoint
		} else {
			// 默认构建endpoint
			endpoint = "http://" + worker.IP + ":" + strconv.Itoa(worker.Port)
		}
		proxyId = worker.WorkerID // 向后兼容，使用WorkerID作为ProxyID
	} else {
		// 降级到旧的Proxy选择逻辑（向后兼容）
		prefer := prof.PreferProxy
		if req.PreferProxy != "" {
			prefer = req.PreferProxy
		}
		endpoint, proxyId = domain.PickProxy(l.svcCtx, prefer)
	}

	// expiry
	ttl := l.svcCtx.Config.Ops.SessionTokenTTL
	if ttl <= 0 {
		ttl = 1800
	}
	exp := time.Now().Add(time.Duration(ttl) * time.Second).Unix()

	// sign jwt
	secret := []byte("")
	if l.svcCtx.Config.Middleware.Auth != nil && l.svcCtx.Config.Middleware.Auth.AccessSecret != "" {
		secret = []byte(l.svcCtx.Config.Middleware.Auth.AccessSecret)
	}
	claims := jwt.MapClaims{"sessionId": newSessionID(), "protocol": req.Protocol, "ciId": req.CiId, "exp": exp}
	if l.svcCtx.ContextManager != nil {
		if tid := l.svcCtx.ContextManager.GetTenantID(l.ctx); tid != "" {
			claims["tenantId"] = tid
		}
		if uid := l.svcCtx.ContextManager.GetUserID(l.ctx); uid != "" {
			claims["userId"] = uid
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}

	// build ws url and signature
	wsPath := "/api/" + req.Protocol + "/websocket"
	if req.Protocol == "ssh" {
		wsPath = "/api/ssh/websocket"
	}
	params := domain.HandshakeParams(prof, req.Protocol)
	// merge extra params/options
	if req.Params != nil {
		for k, v := range req.Params {
			params[k] = v
		}
	}
	if req.Options.Port > 0 {
		params["port"] = strconv.Itoa(req.Options.Port)
	}
	if req.Protocol == "rdp" || req.Protocol == "vnc" {
		if req.Options.Width > 0 {
			params["width"] = strconv.Itoa(req.Options.Width)
		}
		if req.Options.Height > 0 {
			params["height"] = strconv.Itoa(req.Options.Height)
		}
		if req.Options.Depth > 0 {
			params["color_depth"] = strconv.Itoa(req.Options.Depth)
		}
	}
	if req.Options.Locale != "" {
		params["locale"] = req.Options.Locale
	}
	if req.Options.Timezone != "" {
		params["timezone"] = req.Options.Timezone
	}
	params["sig"] = signParams(claims["sessionId"].(string), req.Protocol, req.CiId, params, secret)
	finalURL := joinURL(endpoint, wsPath, params)

	// record session in memory
	if l.svcCtx.SessionStore != nil {
		tenant, user := "", ""
		if l.svcCtx.ContextManager != nil {
			tenant = l.svcCtx.ContextManager.GetTenantID(l.ctx)
			user = l.svcCtx.ContextManager.GetUserID(l.ctx)
		}
		l.svcCtx.SessionStore.Start(&svc.Session{
			ID:        claims["sessionId"].(string),
			TenantId:  tenant,
			UserId:    user,
			CiId:      req.CiId,
			Protocol:  req.Protocol,
			ProxyId:   proxyId,
			Endpoint:  endpoint,
			CreatedAt: time.Now().Unix(),
			ExpiresAt: exp,
			Status:    "active",
		})
	}

	// persist via RPC if available
	if l.svcCtx.OpsClient != nil {
		sid := claims["sessionId"].(string)
		var user string
		if l.svcCtx.ContextManager != nil {
			user = l.svcCtx.ContextManager.GetUserID(l.ctx)
		}

		// Phase 3: 包含Worker信息
		sessionInfo := &ops.SessionInfo{
			SessionId: &sid,
			UserId:    &user,
			CiId:      &req.CiId,
			Protocol:  &req.Protocol,
			ProxyId:   &proxyId,
			Endpoint:  &endpoint,
			ExpiresAt: &exp,
			StatusStr: strPtr("active"),
		}

		// 添加Worker绑定信息（如果使用了WorkerClient）
		if workerID != "" {
			sessionInfo.WorkerId = &workerID
			sessionInfo.WorkerIp = &workerIP
			sessionInfo.WorkerPort = int32Ptr(int32(workerPort))
		}

		_, _ = l.svcCtx.OpsClient.CreateSession(l.ctx, sessionInfo)
	}

	return &types.CreateSessionResp{SessionId: claims["sessionId"].(string), ProxyId: proxyId, WsUrl: finalURL, Token: signed, ExpiresAt: exp, Handshake: params}, nil
}

func newSessionID() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }

func joinURL(base string, path string, params map[string]string) string {
	if base == "" {
		return path
	}
	u, err := url.Parse(base)
	if err != nil {
		return base + path
	}
	u.Path = u.ResolveReference(&url.URL{Path: path}).Path
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func signParams(sessionId, protocol, ciId string, params map[string]string, secret []byte) string {
	// canonicalize
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	// simple bubble sort to avoid extra import
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	canon := ""
	first := true
	for _, k := range keys {
		if k == "sig" {
			continue
		}
		if !first {
			canon += "&"
		} else {
			first = false
		}
		canon += k + "=" + params[k]
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(sessionId + "|" + protocol + "|" + ciId + "|" + canon))
	return hex.EncodeToString(mac.Sum(nil))
}

type badRequest struct{ msg string }

func (e badRequest) Error() string { return e.msg }
func errBadRequest(m string) error { return badRequest{msg: m} }

// IsBadRequest reports whether the error is a client-side validation error.
func IsBadRequest(err error) bool { _, ok := err.(badRequest); return ok }

func strPtr(s string) *string   { return &s }
func int32Ptr(i int32) *int32   { return &i }
