package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"hash/crc32"
	"math/rand"
	"sync/atomic"

	ops "github.com/coder-lulu/newbee-ops-rpc/types/ops"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyPickLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var (
	roundRobinCounter uint64 // 轮询计数器
)

// Proxy选择（需要JWT）
func NewProxyPickLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyPickLogic {
	return &ProxyPickLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyPickLogic) ProxyPick(req *types.ProxyPickReq) (resp *types.ProxyPickResp, err error) {
	// 1. 获取可用的 Proxy 列表
	proxyList, err := l.svcCtx.OpsClient.GetProxyList(l.ctx, &ops.ProxyListReq{
		Page:        1,
		PageSize:    1000,
		ProxyStatus: pointy("online"),
		Status:      pointy(uint32(1)), // 启用状态
	})
	if err != nil {
		l.Logger.Errorw("Failed to get proxy list", logx.Field("error", err))
		return nil, err
	}

	// 2. 过滤符合条件的 Proxy
	candidates := l.filterProxies(proxyList.Data, req)
	if len(candidates) == 0 {
		return nil, errors.New("no available proxy found")
	}

	// 3. 根据策略选择 Proxy
	var selected *ops.ProxyInfo
	switch req.Strategy {
	case "least_connections":
		selected = l.selectLeastConnections(candidates)
	case "round_robin":
		selected = l.selectRoundRobin(candidates)
	case "weighted":
		selected = l.selectWeighted(candidates)
	case "consistent_hash":
		selected = l.selectConsistentHash(candidates, req.SessionID)
	case "random":
		selected = l.selectRandom(candidates)
	default:
		// 默认使用最少连接数策略
		selected = l.selectLeastConnections(candidates)
	}

	if selected == nil {
		return nil, errors.New("failed to select proxy")
	}

	// 4. 构造响应
	// 解析endpoints JSON字符串为map
	endpoints := make(map[string]string)
	if selected.Endpoints != nil && *selected.Endpoints != "" {
		if err := json.Unmarshal([]byte(*selected.Endpoints), &endpoints); err != nil {
			l.Logger.Errorw("Failed to parse endpoints JSON",
				logx.Field("proxy_id", selected.ProxyId),
				logx.Field("endpoints_json", *selected.Endpoints),
				logx.Field("error", err))
		}
	}

	// 解引用指针字段，处理nil情况
	proxyID := ""
	if selected.ProxyId != nil {
		proxyID = *selected.ProxyId
	}
	name := ""
	if selected.Name != nil {
		name = *selected.Name
	}
	ip := ""
	if selected.Ip != nil {
		ip = *selected.Ip
	}
	port := 0
	if selected.Port != nil {
		port = int(*selected.Port)
	}

	l.Logger.Infow("Proxy selected",
		logx.Field("proxy_id", proxyID),
		logx.Field("name", name),
		logx.Field("strategy", req.Strategy))

	return &types.ProxyPickResp{
		Code: 0,
		Msg:  "success",
		Data: types.ProxyPickData{
			ProxyID:   proxyID,
			Name:      name,
			IP:        ip,
			Port:      port,
			Endpoints: endpoints,
		},
	}, nil
}

// filterProxies 过滤符合条件的 Proxy
func (l *ProxyPickLogic) filterProxies(proxies []*ops.ProxyInfo, req *types.ProxyPickReq) []*ops.ProxyInfo {
	result := make([]*ops.ProxyInfo, 0)
	for _, p := range proxies {
		// 过滤排除列表
		if p.ProxyId != nil && containsString(req.ExcludeProxyIDs, *p.ProxyId) {
			continue
		}

		// 过滤能力要求
		if len(req.RequiredCapabilities) > 0 {
			// Capabilities是JSON字符串，需要解析
			var caps []string
			if p.Capabilities != nil && *p.Capabilities != "" {
				json.Unmarshal([]byte(*p.Capabilities), &caps)
			}
			if !hasCapabilitiesSlice(caps, req.RequiredCapabilities) {
				continue
			}
		}

		// 过滤优先区域
		if req.PreferredRegion != "" && (p.Region == nil || *p.Region != req.PreferredRegion) {
			continue
		}

		// 过滤优先可用区
		if req.PreferredZone != "" && (p.Zone == nil || *p.Zone != req.PreferredZone) {
			continue
		}

		// 过滤必须的标签
		if len(req.RequiredTags) > 0 {
			// Tags是JSON字符串，需要解析
			var tags []string
			if p.Tags != nil && *p.Tags != "" {
				json.Unmarshal([]byte(*p.Tags), &tags)
			}
			if !hasAllTagsSlice(tags, req.RequiredTags) {
				continue
			}
		}

		// 过滤最小权重
		if req.MinWeight > 0 && (p.Weight == nil || int(*p.Weight) < req.MinWeight) {
			continue
		}

		result = append(result, p)
	}
	return result
}

// selectLeastConnections 选择连接数最少的 Proxy
func (l *ProxyPickLogic) selectLeastConnections(candidates []*ops.ProxyInfo) *ops.ProxyInfo {
	if len(candidates) == 0 {
		return nil
	}

	var selected *ops.ProxyInfo
	minSessions := int64(999999)
	for _, p := range candidates {
		sessions := int64(0)
		if p.ActiveSessions != nil {
			sessions = *p.ActiveSessions
		}
		if sessions < minSessions {
			minSessions = sessions
			selected = p
		}
	}
	return selected
}

// selectRoundRobin 轮询选择 Proxy
func (l *ProxyPickLogic) selectRoundRobin(candidates []*ops.ProxyInfo) *ops.ProxyInfo {
	if len(candidates) == 0 {
		return nil
	}

	index := atomic.AddUint64(&roundRobinCounter, 1) % uint64(len(candidates))
	return candidates[index]
}

// selectWeighted 加权轮询选择 Proxy
func (l *ProxyPickLogic) selectWeighted(candidates []*ops.ProxyInfo) *ops.ProxyInfo {
	if len(candidates) == 0 {
		return nil
	}

	// 计算总权重
	var totalWeight int64
	for _, p := range candidates {
		if p.Weight != nil {
			totalWeight += *p.Weight
		}
	}

	if totalWeight == 0 {
		// 如果总权重为0，使用轮询
		return l.selectRoundRobin(candidates)
	}

	// 生成随机数
	r := rand.Int63n(totalWeight)
	var sum int64
	for _, p := range candidates {
		weight := int64(0)
		if p.Weight != nil {
			weight = *p.Weight
		}
		sum += weight
		if r < sum {
			return p
		}
	}

	return candidates[0]
}

// selectConsistentHash 一致性哈希选择 Proxy
func (l *ProxyPickLogic) selectConsistentHash(candidates []*ops.ProxyInfo, sessionID string) *ops.ProxyInfo {
	if len(candidates) == 0 {
		return nil
	}

	if sessionID == "" {
		// 如果没有 sessionID，使用轮询
		return l.selectRoundRobin(candidates)
	}

	// 计算哈希值
	hash := crc32.ChecksumIEEE([]byte(sessionID))
	index := int(hash) % len(candidates)
	return candidates[index]
}

// selectRandom 随机选择 Proxy
func (l *ProxyPickLogic) selectRandom(candidates []*ops.ProxyInfo) *ops.ProxyInfo {
	if len(candidates) == 0 {
		return nil
	}

	index := rand.Intn(len(candidates))
	return candidates[index]
}

// containsString 检查字符串是否在列表中
func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// hasCapabilitiesSlice 检查是否包含所有必需的能力
func hasCapabilitiesSlice(capabilities []string, required []string) bool {
	capSet := make(map[string]bool)
	for _, cap := range capabilities {
		capSet[cap] = true
	}
	for _, req := range required {
		if !capSet[req] {
			return false
		}
	}
	return true
}

// hasAllTagsSlice 检查是否包含所有必需的标签
func hasAllTagsSlice(tags []string, required []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}
	for _, req := range required {
		if !tagSet[req] {
			return false
		}
	}
	return true
}

// hasCapabilities 检查是否包含所有必须的能力
func hasCapabilities(proxyCapabilities []string, required []string) bool {
	capMap := make(map[string]bool)
	for _, cap := range proxyCapabilities {
		capMap[cap] = true
	}

	for _, req := range required {
		if !capMap[req] {
			return false
		}
	}
	return true
}

// hasAllTags 检查是否包含所有必须的标签
func hasAllTags(proxyTags []string, required []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range proxyTags {
		tagMap[tag] = true
	}

	for _, req := range required {
		if !tagMap[req] {
			return false
		}
	}
	return true
}

// pointy helper function
