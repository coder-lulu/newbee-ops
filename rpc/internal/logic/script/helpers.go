package script

import (
	"encoding/json"
	"fmt"
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

// newInvalidArgumentError creates a simple error for invalid arguments
func newInvalidArgumentError(field string) error {
	return fmt.Errorf("invalid argument: %s is required", field)
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
