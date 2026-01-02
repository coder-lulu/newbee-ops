package script_version

import (
	"encoding/json"
)

// mustMarshalJSON marshals a value to JSON string
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

// unmarshalJSONMap unmarshals JSON string to map[string]interface{}
func unmarshalJSONMap(jsonStr *string) map[string]interface{} {
	if jsonStr == nil || *jsonStr == "" {
		return map[string]interface{}{}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}
