package script

import (
	"encoding/json"
)

// pointy helper function to convert values to pointers
func pointy[T any](v T) *T {
	return &v
}

// getValue helper function to safely dereference pointers
func getValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

// marshalJSONStringSlice marshals []string to JSON string pointer
func marshalJSONStringSlice(slice []string) *string {
	if len(slice) == 0 {
		return nil
	}
	data, err := json.Marshal(slice)
	if err != nil {
		return nil
	}
	result := string(data)
	return &result
}

// marshalJSONMap marshals map[string]string to JSON string pointer
func marshalJSONMap(m map[string]string) *string {
	if len(m) == 0 {
		return nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	result := string(data)
	return &result
}

// unmarshalJSONStringSlice unmarshals JSON string pointer to []string
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

// unmarshalJSONMap unmarshals JSON string pointer to map[string]string
func unmarshalJSONMap(jsonStr *string) map[string]string {
	if jsonStr == nil || *jsonStr == "" {
		return map[string]string{}
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		return map[string]string{}
	}
	return result
}
