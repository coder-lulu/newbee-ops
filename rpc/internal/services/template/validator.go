package template

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"github.com/zeromicro/go-zero/core/logx"
)

// ParameterValidator 参数验证器接口
type ParameterValidator interface {
	// Validate 验证参数
	Validate(schemaJSON string, parameters map[string]interface{}) error
	// ValidateWithDetail 验证参数并返回详细错误
	ValidateWithDetail(schemaJSON string, parameters map[string]interface{}) (bool, []string, error)
}

// parameterValidatorImpl 参数验证器实现
type parameterValidatorImpl struct {
	logger logx.Logger
}

// NewParameterValidator 创建参数验证器实例
func NewParameterValidator() ParameterValidator {
	return &parameterValidatorImpl{
		logger: logx.WithContext(nil),
	}
}

// Validate 验证参数
func (v *parameterValidatorImpl) Validate(schemaJSON string, parameters map[string]interface{}) error {
	valid, errors, err := v.ValidateWithDetail(schemaJSON, parameters)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("parameter validation failed: %v", errors)
	}

	return nil
}

// ValidateWithDetail 验证参数并返回详细错误
func (v *parameterValidatorImpl) ValidateWithDetail(schemaJSON string, parameters map[string]interface{}) (bool, []string, error) {
	// 如果没有schema，跳过验证
	if schemaJSON == "" || schemaJSON == "{}" {
		v.logger.Debug("No schema provided, skip validation")
		return true, nil, nil
	}

	// 加载Schema
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)

	// 转换参数为JSON
	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	// 加载参数数据
	documentLoader := gojsonschema.NewStringLoader(string(parametersJSON))

	// 执行验证
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, nil, fmt.Errorf("validation error: %w", err)
	}

	// 检查验证结果
	if result.Valid() {
		v.logger.Debug("Parameters validation passed")
		return true, nil, nil
	}

	// 收集错误信息
	var errorMessages []string
	for _, err := range result.Errors() {
		errorMessages = append(errorMessages, err.String())
	}

	v.logger.Debugw("Parameters validation failed", logx.Field("errors", errorMessages))
	return false, errorMessages, nil
}

// ValidateSchema 验证Schema本身是否合法
func (v *parameterValidatorImpl) ValidateSchema(schemaJSON string) error {
	// 尝试解析Schema
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		return fmt.Errorf("invalid JSON schema: %w", err)
	}

	// 尝试加载Schema（验证格式）
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	_, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("invalid JSON schema format: %w", err)
	}

	return nil
}

// ParameterDefinition 参数定义
type ParameterDefinition struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // string, number, boolean, array, object
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	MinLength   *int        `json:"minLength,omitempty"`
	MaxLength   *int        `json:"maxLength,omitempty"`
	Minimum     *float64    `json:"minimum,omitempty"`
	Maximum     *float64    `json:"maximum,omitempty"`
}

// BuildJSONSchema 从参数定义构建JSON Schema
func BuildJSONSchema(parameters []ParameterDefinition) (string, error) {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type":    "object",
	}

	properties := make(map[string]interface{})
	var required []string

	for _, param := range parameters {
		// 构建参数属性
		property := map[string]interface{}{
			"type": param.Type,
		}

		if param.Description != "" {
			property["description"] = param.Description
		}

		if param.Default != nil {
			property["default"] = param.Default
		}

		if len(param.Enum) > 0 {
			property["enum"] = param.Enum
		}

		if param.Pattern != "" {
			property["pattern"] = param.Pattern
		}

		if param.MinLength != nil {
			property["minLength"] = *param.MinLength
		}

		if param.MaxLength != nil {
			property["maxLength"] = *param.MaxLength
		}

		if param.Minimum != nil {
			property["minimum"] = *param.Minimum
		}

		if param.Maximum != nil {
			property["maximum"] = *param.Maximum
		}

		properties[param.Name] = property

		// 添加到必填列表
		if param.Required {
			required = append(required, param.Name)
		}
	}

	schema["properties"] = properties

	if len(required) > 0 {
		schema["required"] = required
	}

	// 转换为JSON字符串
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(schemaBytes), nil
}

// ExtractDefaults 从参数定义中提取默认值
func ExtractDefaults(parameters []ParameterDefinition) map[string]interface{} {
	defaults := make(map[string]interface{})

	for _, param := range parameters {
		if param.Default != nil {
			defaults[param.Name] = param.Default
		}
	}

	return defaults
}

// MergeWithDefaults 合并参数和默认值
func MergeWithDefaults(parameters map[string]interface{}, defaults map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 先复制默认值
	for k, v := range defaults {
		result[k] = v
	}

	// 用实际参数覆盖
	for k, v := range parameters {
		result[k] = v
	}

	return result
}
