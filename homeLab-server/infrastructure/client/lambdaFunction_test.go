package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLambdaFnConfig_ValidConfig(t *testing.T) {
	// Test valid configuration
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action":     "turn_on",
			"ip_address": "192.168.1.10",
		},
		"base_url": "http://localhost:8080",
		"token":    "test-token",
		"headers": map[string]string{
			"X-Custom-Header": "test-value",
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/test-function", result.FunctionName)
	assert.Equal(t, "http://localhost:8080", result.BaseURL)
	assert.Equal(t, "test-token", result.Token)
	assert.Equal(t, "turn_on", result.Args["action"])
	assert.Equal(t, "192.168.1.10", result.Args["ip_address"])
	assert.Equal(t, "test-value", result.Headers["X-Custom-Header"])
}

func TestNewLambdaFnConfig_MissingFunctionName(t *testing.T) {
	// Test configuration without function name
	config := map[string]interface{}{
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid or missing function name")
}

func TestNewLambdaFnConfig_InvalidFunctionName(t *testing.T) {
	// Test configuration with invalid function name type
	config := map[string]interface{}{
		"fn_name": 123, // Should be string
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid or missing function name")
}

func TestNewLambdaFnConfig_MinimalConfig(t *testing.T) {
	// Test minimal configuration with only required fields
	config := map[string]interface{}{
		"fn_name": "function/minimal-function",
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/minimal-function", result.FunctionName)
	assert.Empty(t, result.BaseURL)
	assert.Empty(t, result.Token)
	assert.Nil(t, result.Args)
	assert.Empty(t, result.Headers)
}

func TestNewLambdaFnConfig_WithArgsOnly(t *testing.T) {
	// Test configuration with args only
	config := map[string]interface{}{
		"fn_name": "function/args-function",
		"args": map[string]interface{}{
			"action":   "discover_devices",
			"timeout":  30,
			"retry":    true,
			"ip_range": "192.168.1.0/24",
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/args-function", result.FunctionName)
	assert.Equal(t, "discover_devices", result.Args["action"])
	assert.Equal(t, float64(30), result.Args["timeout"]) // JSON unmarshaling converts numbers to float64
	assert.Equal(t, true, result.Args["retry"])
	assert.Equal(t, "192.168.1.0/24", result.Args["ip_range"])
}

func TestNewLambdaFnConfig_WithHeadersOnly(t *testing.T) {
	// Test configuration with headers only
	config := map[string]interface{}{
		"fn_name": "function/headers-function",
		"headers": map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
			"X-API-Key":     "api-key-123",
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/headers-function", result.FunctionName)
	assert.Equal(t, "Bearer token123", result.Headers["Authorization"])
	assert.Equal(t, "application/json", result.Headers["Content-Type"])
	assert.Equal(t, "api-key-123", result.Headers["X-API-Key"])
}

func TestNewLambdaFnConfigUnmarshal_ValidJSON(t *testing.T) {
	// Test unmarshaling from JSON string
	jsonConfig := `{
		"fn_name": "function/json-function",
		"args": {
			"action": "get_system_info",
			"timeout": 60
		},
		"base_url": "https://api.example.com",
		"token": "json-token"
	}`

	result, err := NewLambdaFnConfigUnmarshal(jsonConfig)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/json-function", result.FunctionName)
	assert.Equal(t, "get_system_info", result.Args["action"])
	assert.Equal(t, float64(60), result.Args["timeout"])
	assert.Equal(t, "https://api.example.com", result.BaseURL)
	assert.Equal(t, "json-token", result.Token)
}

func TestNewLambdaFnConfigUnmarshal_InvalidJSON(t *testing.T) {
	// Test unmarshaling from invalid JSON string
	invalidJSON := `{
		"fn_name": "function/invalid-function",
		"args": {
			"action": "test_action",
		}
	}` // Missing closing brace

	result, err := NewLambdaFnConfigUnmarshal(invalidJSON)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNewLambdaFnConfig_ComplexArgs(t *testing.T) {
	// Test configuration with complex nested arguments
	config := map[string]interface{}{
		"fn_name": "function/complex-function",
		"args": map[string]interface{}{
			"action": "deploy_service",
			"service": map[string]interface{}{
				"name":    "satisfactory",
				"version": "1.0.0",
				"config": map[string]interface{}{
					"port":        7777,
					"max_players": 8,
					"auto_save":   true,
				},
			},
			"targets": []string{"node-1", "node-2"},
			"timeout": 300,
		},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/complex-function", result.FunctionName)

	// Check nested service config
	service := result.Args["service"].(map[string]interface{})
	assert.Equal(t, "satisfactory", service["name"])
	assert.Equal(t, "1.0.0", service["version"])

	serviceConfig := service["config"].(map[string]interface{})
	assert.Equal(t, float64(7777), serviceConfig["port"])
	assert.Equal(t, float64(8), serviceConfig["max_players"])
	assert.Equal(t, true, serviceConfig["auto_save"])

	// Check targets array
	targets := result.Args["targets"].([]interface{})
	assert.Equal(t, "node-1", targets[0])
	assert.Equal(t, "node-2", targets[1])

	assert.Equal(t, float64(300), result.Args["timeout"])
}

func TestNewLambdaFnConfig_RealWorldExample(t *testing.T) {
	// Test with a real-world example from your pipeline config
	config := map[string]interface{}{
		"fn_name": "function/network-controller",
		"args": map[string]interface{}{
			"action":     "turn_on",
			"ip_address": "192.168.1.10",
		},
		"success_condition_key":   "device.status",
		"success_condition_value": "on",
		"interval":                5,
		"retry":                   3,
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/network-controller", result.FunctionName)
	assert.Equal(t, "turn_on", result.Args["action"])
	assert.Equal(t, "192.168.1.10", result.Args["ip_address"])
	assert.Equal(t, "device.status", config["success_condition_key"])
	assert.Equal(t, "on", config["success_condition_value"])
	assert.Equal(t, float64(5), config["interval"])
	assert.Equal(t, float64(3), config["retry"])
}

func TestNewLambdaFnConfig_EmptyArgs(t *testing.T) {
	// Test configuration with empty args
	config := map[string]interface{}{
		"fn_name": "function/empty-args-function",
		"args":    map[string]interface{}{},
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/empty-args-function", result.FunctionName)
	assert.NotNil(t, result.Args)
	assert.Empty(t, result.Args)
}

func TestNewLambdaFnConfig_NilArgs(t *testing.T) {
	// Test configuration with nil args
	config := map[string]interface{}{
		"fn_name": "function/nil-args-function",
		"args":    nil,
	}

	result, err := NewLambdaFnConfig(config)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "function/nil-args-function", result.FunctionName)
	assert.Nil(t, result.Args)
}

// Test helper function to create test config
func createTestLambdaConfig(fnName string, args map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"fn_name": fnName,
		"args":    args,
	}
}

// Test helper function to create test config with all optional fields
func createFullTestLambdaConfig(fnName string, args map[string]interface{}, baseURL, token string, headers map[string]string) map[string]interface{} {
	config := map[string]interface{}{
		"fn_name": fnName,
		"args":    args,
	}

	if baseURL != "" {
		config["base_url"] = baseURL
	}

	if token != "" {
		config["token"] = token
	}

	if headers != nil {
		config["headers"] = headers
	}

	return config
}
