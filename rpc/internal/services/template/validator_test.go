package template

import (
	"testing"
)

func TestParameterValidator_Validate(t *testing.T) {
	validator := NewParameterValidator()

	// JSON Schema示例：需要name（字符串）和age（数字）
	schema := `{
		"type": "object",
		"properties": {
			"name": {
				"type": "string",
				"minLength": 1
			},
			"age": {
				"type": "number",
				"minimum": 0,
				"maximum": 150
			},
			"email": {
				"type": "string",
				"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
			}
		},
		"required": ["name", "age"]
	}`

	tests := []struct {
		name        string
		schema      string
		parameters  map[string]interface{}
		expectValid bool
	}{
		{
			name:   "Valid parameters",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John Doe",
				"age":  30.0,
			},
			expectValid: true,
		},
		{
			name:   "Valid with email",
			schema: schema,
			parameters: map[string]interface{}{
				"name":  "Jane Doe",
				"age":   25.0,
				"email": "jane@example.com",
			},
			expectValid: true,
		},
		{
			name:   "Missing required field",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John",
			},
			expectValid: false,
		},
		{
			name:   "Invalid type - age is string",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John",
				"age":  "thirty",
			},
			expectValid: false,
		},
		{
			name:   "Age out of range",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John",
				"age":  200.0,
			},
			expectValid: false,
		},
		{
			name:   "Invalid email format",
			schema: schema,
			parameters: map[string]interface{}{
				"name":  "John",
				"age":   30.0,
				"email": "invalid-email",
			},
			expectValid: false,
		},
		{
			name:        "Empty schema - skip validation",
			schema:      "",
			parameters:  map[string]interface{}{"any": "value"},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.schema, tt.parameters)

			if tt.expectValid && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}

			if !tt.expectValid && err == nil {
				t.Errorf("Expected validation to fail, but it passed")
			}
		})
	}
}

func TestParameterValidator_ValidateWithDetail(t *testing.T) {
	validator := NewParameterValidator()

	schema := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "number", "minimum": 0}
		},
		"required": ["name"]
	}`

	tests := []struct {
		name         string
		schema       string
		parameters   map[string]interface{}
		expectValid  bool
		expectErrors bool
	}{
		{
			name:   "Valid parameters",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John",
				"age":  30.0,
			},
			expectValid:  true,
			expectErrors: false,
		},
		{
			name:   "Missing required field",
			schema: schema,
			parameters: map[string]interface{}{
				"age": 30.0,
			},
			expectValid:  false,
			expectErrors: true,
		},
		{
			name:   "Invalid type",
			schema: schema,
			parameters: map[string]interface{}{
				"name": "John",
				"age":  "thirty",
			},
			expectValid:  false,
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, errors, err := validator.ValidateWithDetail(tt.schema, tt.parameters)

			if err != nil {
				t.Fatalf("Unexpected validation error: %v", err)
			}

			if valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, valid)
			}

			if tt.expectErrors && len(errors) == 0 {
				t.Errorf("Expected error messages, but got none")
			}

			if !tt.expectErrors && len(errors) > 0 {
				t.Errorf("Expected no error messages, but got: %v", errors)
			}
		})
	}
}

func TestBuildJSONSchema(t *testing.T) {
	params := []ParameterDefinition{
		{
			Name:        "hostname",
			Type:        "string",
			Required:    true,
			Description: "Target hostname",
			Pattern:     "^[a-zA-Z0-9.-]+$",
		},
		{
			Name:        "port",
			Type:        "number",
			Required:    false,
			Default:     22.0,
			Description: "Port number",
			Minimum:     getFloat64Pointer(1),
			Maximum:     getFloat64Pointer(65535),
		},
		{
			Name:        "protocol",
			Type:        "string",
			Required:    true,
			Description: "Protocol to use",
			Enum:        []string{"ssh", "telnet", "rdp"},
		},
	}

	schema, err := BuildJSONSchema(params)
	if err != nil {
		t.Fatalf("BuildJSONSchema failed: %v", err)
	}

	if schema == "" {
		t.Errorf("Expected non-empty schema")
	}

	// 验证生成的schema是合法的
	validator := NewParameterValidator()
	validParams := map[string]interface{}{
		"hostname": "example.com",
		"port":     22.0,
		"protocol": "ssh",
	}

	if err := validator.Validate(schema, validParams); err != nil {
		t.Errorf("Generated schema validation failed: %v", err)
	}
}

func TestExtractDefaults(t *testing.T) {
	params := []ParameterDefinition{
		{
			Name:     "hostname",
			Type:     "string",
			Required: true,
		},
		{
			Name:     "port",
			Type:     "number",
			Required: false,
			Default:  22.0,
		},
		{
			Name:     "timeout",
			Type:     "number",
			Required: false,
			Default:  30.0,
		},
	}

	defaults := ExtractDefaults(params)

	if len(defaults) != 2 {
		t.Errorf("Expected 2 defaults, got %d", len(defaults))
	}

	if defaults["port"] != 22.0 {
		t.Errorf("Expected port=22, got %v", defaults["port"])
	}

	if defaults["timeout"] != 30.0 {
		t.Errorf("Expected timeout=30, got %v", defaults["timeout"])
	}
}

func TestMergeWithDefaults(t *testing.T) {
	defaults := map[string]interface{}{
		"port":    22.0,
		"timeout": 30.0,
		"retries": 3.0,
	}

	parameters := map[string]interface{}{
		"hostname": "example.com",
		"port":     2222.0, // Override default
	}

	result := MergeWithDefaults(parameters, defaults)

	// Should have 4 keys
	if len(result) != 4 {
		t.Errorf("Expected 4 keys, got %d", len(result))
	}

	// hostname from parameters
	if result["hostname"] != "example.com" {
		t.Errorf("Expected hostname=example.com")
	}

	// port overridden
	if result["port"] != 2222.0 {
		t.Errorf("Expected port=2222, got %v", result["port"])
	}

	// timeout from defaults
	if result["timeout"] != 30.0 {
		t.Errorf("Expected timeout=30, got %v", result["timeout"])
	}

	// retries from defaults
	if result["retries"] != 3.0 {
		t.Errorf("Expected retries=3, got %v", result["retries"])
	}
}

// Helper function to create float64 pointer
func getFloat64Pointer(v float64) *float64 {
	return &v
}
