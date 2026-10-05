package domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/coder-lulu/newbee-ops-api/internal/svc"
	"github.com/coder-lulu/newbee-ops-api/internal/types"
	"github.com/coder-lulu/newbee-ops-rpc/types/ops"
	"google.golang.org/protobuf/types/known/structpb"
)

// LoadProfile 从 RPC（优先）或内存加载画像
func LoadProfile(ctx context.Context, svcCtx *svc.ServiceContext, ciId string) (types.AccessProfile, bool, error) {
	if svcCtx.OpsClient != nil {
		// 使用列表过滤以按 ciId 获取（与当前 RPC 定义对齐）
		res, err := svcCtx.OpsClient.GetAccessProfileList(ctx, &ops.AccessProfileListReq{Page: 1, PageSize: 1, CiId: &ciId})
		if err != nil {
			return types.AccessProfile{}, false, err
		}
		if res != nil && len(res.Data) > 0 {
			item := res.Data[0]
			ap := types.AccessProfile{CiId: item.GetCiId(), CredentialRef: item.GetCredentialRef(), PreferProxy: item.GetPreferProxy()}
			if s := item.GetCapabilities(); s != nil {
				ap.Capabilities = unpackStringList(s)
			}
			if s := item.GetPorts(); s != nil {
				ap.Ports = unpackIntMap(s)
			}
			if s := item.GetJumpChain(); s != nil {
				ap.JumpChain = unpackStringList(s)
			}
			if s := item.GetTags(); s != nil {
				ap.Tags = unpackStringMap(s)
			}
			return ap, true, nil
		}
		return types.AccessProfile{}, false, nil
	}
	if v, ok := svcCtx.ProfileStore.Get(ciId); ok {
		// svc.AccessProfile is now compatible with types.AccessProfile
		// Manual conversion if needed, but since we updated svc.AccessProfile, we can try direct assignment if compatible,
		// or copy fields. Since types.AccessProfile is defined in types package, and svc.AccessProfile in svc, they are different types.
		ap := types.AccessProfile{
			CiId:          v.CiId,
			Capabilities:  v.Capabilities,
			Ports:         v.Ports,
			CredentialRef: v.CredentialRef,
			PreferProxy:   v.PreferProxy,
			JumpChain:     v.JumpChain,
			Tags:          v.Tags,
		}
		return ap, true, nil
	}
	return types.AccessProfile{}, false, nil
}

func unpackStringList(s *structpb.Struct) []string {
	m := s.AsMap()
	if list, ok := m["list"].([]interface{}); ok {
		out := make([]string, 0, len(list))
		for _, v := range list {
			if str, ok := v.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}

func unpackIntMap(s *structpb.Struct) map[string]int {
	m := s.AsMap()
	out := make(map[string]int, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case float64:
			out[k] = int(val)
		case int:
			out[k] = val
		case int64:
			out[k] = int(val)
		}
	}
	return out
}

func unpackStringMap(s *structpb.Struct) map[string]string {
	m := s.AsMap()
	out := make(map[string]string, len(m))
	for k, v := range m {
		if str, ok := v.(string); ok {
			out[k] = str
		} else {
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out
}

// AllowProtocol 校验协议是否被画像允许（空 capabilities 视为放行）
func AllowProtocol(p types.AccessProfile, protocol string) bool {
	if len(p.Capabilities) == 0 {
		return true
	}
	for _, c := range p.Capabilities {
		if c == protocol {
			return true
		}
	}
	return false
}

// HandshakeParams 计算握手参数（例如端口）
func HandshakeParams(p types.AccessProfile, protocol string) map[string]string {
	out := map[string]string{}
	if p.Ports != nil {
		if port, ok := p.Ports[protocol]; ok && port > 0 {
			out["port"] = strconv.Itoa(port)
		}
	}
	// additional parameters from tags
	if protocol == "rdp" {
		if w, ok := p.Tags["rdp_width"]; ok {
			out["width"] = w
		}
		if h, ok := p.Tags["rdp_height"]; ok {
			out["height"] = h
		}
		if d, ok := p.Tags["rdp_dpi"]; ok {
			out["dpi"] = d
		}
		if c, ok := p.Tags["rdp_color_depth"]; ok {
			out["color-depth"] = c
		}
	}
	if protocol == "vnc" {
		if w, ok := p.Tags["vnc_width"]; ok {
			out["width"] = w
		}
		if h, ok := p.Tags["vnc_height"]; ok {
			out["height"] = h
		}
	}
	if protocol == "ssh" {
		if cols, ok := p.Tags["ssh_cols"]; ok {
			out["cols"] = cols
		}
		if rows, ok := p.Tags["ssh_rows"]; ok {
			out["rows"] = rows
		}
	}
	// jump chain passed as json when present
	if len(p.JumpChain) > 0 {
		// minimal json array without import
		s := "["
		for i, j := range p.JumpChain {
			if i > 0 {
				s += ","
			}
			s += "\"" + j + "\""
		}
		s += "]"
		out["jump"] = s
	}
	return out
}

// PickProxy 根据 preferProxy/注册表/默认配置挑选代理
func PickProxy(svcCtx *svc.ServiceContext, preferred string) (endpoint string, proxyId string) {
	endpoint = "http://127.0.0.1:8889"
	proxyId = "proxy-default"
	// preferProxy
	if preferred != "" && svcCtx.ProxyRegistry != nil {
		for _, it := range svcCtx.ProxyRegistry.List() {
			if it.ProxyID == preferred {
				if ep := it.Endpoints["http"]; ep != "" {
					endpoint = ep
				}
				proxyId = it.ProxyID
				return
			}
		}
	}
	// registry pick
	if svcCtx.ProxyRegistry != nil {
		if picked := svcCtx.ProxyRegistry.Pick(); picked != nil {
			if ep := picked.Endpoints["http"]; ep != "" {
				endpoint = ep
			}
			proxyId = picked.ProxyID
			return
		}
	}
	// default endpoint
	if len(svcCtx.Config.Ops.DefaultProxyEndpoints) > 0 {
		endpoint = svcCtx.Config.Ops.DefaultProxyEndpoints[0]
		proxyId = "proxy-001"
		return
	}
	return
}
