package proxy

import (
	"encoding/json"

	"github.com/coder-lulu/newbee-ops-rpc/ent/proxy"
)

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

// unmarshalJSONStringSlice unmarshals JSON string to []string
func unmarshalJSONStringSlice(jsonStr *string) []string {
	if jsonStr == nil || *jsonStr == "" {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		return []string{}
	}
	return result
}

// unmarshalJSONStringMap unmarshals JSON string to map[string]string
func unmarshalJSONStringMap(jsonStr *string) map[string]string {
	if jsonStr == nil || *jsonStr == "" {
		return map[string]string{}
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		return map[string]string{}
	}
	return result
}

// unmarshalJSONInterfaceMap unmarshals JSON string to map[string]interface{}
func unmarshalJSONInterfaceMap(jsonStr *string) map[string]interface{} {
	if jsonStr == nil || *jsonStr == "" {
		return map[string]interface{}{}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}

// stringToProxyStatus converts string to proxy.WorkerStatus enum
func stringToProxyStatus(s *string) proxy.WorkerStatus {
	if s == nil {
		return proxy.WorkerStatusOffline
	}
	switch *s {
	case "online":
		return proxy.WorkerStatusOnline
	case "offline":
		return proxy.WorkerStatusOffline
	default:
		return proxy.WorkerStatusOffline
	}
}

// proxyStatusToString converts proxy.WorkerStatus enum to string
func proxyStatusToString(status proxy.WorkerStatus) string {
	return string(status)
}
