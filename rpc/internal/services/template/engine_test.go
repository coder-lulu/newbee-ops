package template

import (
	"testing"
)

func TestTemplateEngine_Render(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name        string
		template    string
		parameters  map[string]interface{}
		expected    string
		expectError bool
	}{
		{
			name:        "Simple variable substitution",
			template:    "Hello {{.name}}!",
			parameters:  map[string]interface{}{"name": "World"},
			expected:    "Hello World!",
			expectError: false,
		},
		{
			name:        "Default value function",
			template:    "Hello {{.name | default \"Guest\"}}!",
			parameters:  map[string]interface{}{},
			expected:    "Hello Guest!",
			expectError: false,
		},
		{
			name:        "Conditional rendering",
			template:    "{{if .enabled}}Feature is enabled{{else}}Feature is disabled{{end}}",
			parameters:  map[string]interface{}{"enabled": true},
			expected:    "Feature is enabled",
			expectError: false,
		},
		{
			name:        "Loop rendering",
			template:    "{{range .items}}{{.}},{{end}}",
			parameters:  map[string]interface{}{"items": []string{"apple", "banana", "cherry"}},
			expected:    "apple,banana,cherry,",
			expectError: false,
		},
		{
			name:        "Custom function - upper",
			template:    "{{.name | upper}}",
			parameters:  map[string]interface{}{"name": "hello"},
			expected:    "HELLO",
			expectError: false,
		},
		{
			name:        "Custom function - lower",
			template:    "{{.name | lower}}",
			parameters:  map[string]interface{}{"name": "WORLD"},
			expected:    "world",
			expectError: false,
		},
		{
			name:        "Invalid template syntax",
			template:    "{{.name",
			parameters:  map[string]interface{}{"name": "test"},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Render(tt.template, tt.parameters)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateEngine_Validate(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name        string
		template    string
		expectError bool
	}{
		{
			name:        "Valid template",
			template:    "Hello {{.name}}!",
			expectError: false,
		},
		{
			name:        "Invalid template - unclosed tag",
			template:    "Hello {{.name",
			expectError: true,
		},
		{
			name:        "Invalid template - unknown function",
			template:    "Hello {{.name | unknownFunc}}",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.Validate(tt.template)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTemplateEngine_RenderWithCache(t *testing.T) {
	engine := NewTemplateEngine()

	templateContent := "Hello {{.name}}!"
	parameters := map[string]interface{}{"name": "World"}

	// First render - should cache
	result1, err := engine.RenderWithCache("test-template", templateContent, parameters)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Second render - should use cache
	result2, err := engine.RenderWithCache("test-template", templateContent, parameters)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result1 != result2 {
		t.Errorf("Expected same result from cache, got different")
	}

	if result1 != "Hello World!" {
		t.Errorf("Expected %q, got %q", "Hello World!", result1)
	}
}

func TestTemplateEngine_CustomFunctions(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name       string
		template   string
		parameters map[string]interface{}
		expected   string
	}{
		{
			name:       "contains function",
			template:   "{{if contains .text \"hello\"}}Found{{else}}Not found{{end}}",
			parameters: map[string]interface{}{"text": "hello world"},
			expected:   "Found",
		},
		{
			name:       "trim function",
			template:   "{{.text | trim}}",
			parameters: map[string]interface{}{"text": "  hello  "},
			expected:   "hello",
		},
		{
			name:       "hasPrefix function",
			template:   "{{if hasPrefix .text \"hello\"}}True{{else}}False{{end}}",
			parameters: map[string]interface{}{"text": "hello world"},
			expected:   "True",
		},
		{
			name:       "hasSuffix function",
			template:   "{{if hasSuffix .text \"world\"}}True{{else}}False{{end}}",
			parameters: map[string]interface{}{"text": "hello world"},
			expected:   "True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Render(tt.template, tt.parameters)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

