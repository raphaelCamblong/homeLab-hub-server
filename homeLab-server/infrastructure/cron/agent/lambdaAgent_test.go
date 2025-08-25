package agent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

// MockLambdaFunction is a mock implementation of the LambdaFunction interface
type MockLambdaFunction struct {
	mock.Mock
}

func (m *MockLambdaFunction) Invoke(fnCfg client.LambdaFnConfig) ([]byte, error) {
	args := m.Called(fnCfg)
	return args.Get(0).([]byte), args.Error(1)
}

// MockConfig provides a mock configuration for testing
type MockConfig struct {
	Lambda config.LambdaConfig
}

func (m *MockConfig) Get() *config.Config {
	return &config.Config{
		Client: config.ClientConfig{
			Lambda: m.Lambda,
		},
	}
}

func TestNewLambdaAgent(t *testing.T) {
	// Test that NewLambdaAgent creates a valid agent
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-step",
		},
	}

	agent := NewLambdaAgent(step)
	assert.NotNil(t, agent)
	assert.IsType(t, &LambdaAgent{}, agent)

	lambdaAgent := agent.(*LambdaAgent)
	assert.Equal(t, step, lambdaAgent.step)
	assert.NotNil(t, lambdaAgent.lambdaClient)
}

func TestLambdaAgent_Execute_Success(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock response data
	responseData := map[string]interface{}{
		"status": "success",
		"data":   "test data",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client to return success
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Invoking lambda function: test-lambda-step")
	assert.Contains(t, result.Log, "test data")

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_WithSuccessCondition(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock response data that matches success condition
	responseData := map[string]interface{}{
		"device": map[string]interface{}{
			"status": "on",
		},
		"message": "Device turned on successfully",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config with success condition (matching your pipeline config)
	config := map[string]interface{}{
		"fn_name": "function/network-controller",
		"args": map[string]interface{}{
			"action":     "turn_on",
			"ip_address": "192.168.1.10",
		},
		"success_condition_key":   "device.status",
		"success_condition_value": "on",
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Device turned on successfully")

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_WithFailedSuccessCondition(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock response data that doesn't match success condition
	responseData := map[string]interface{}{
		"device": map[string]interface{}{
			"status": "off", // Different from expected "on"
		},
		"message": "Device is still off",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config with success condition
	config := map[string]interface{}{
		"fn_name": "function/network-controller",
		"args": map[string]interface{}{
			"action":     "turn_on",
			"ip_address": "192.168.1.10",
		},
		"success_condition_key":   "device.status",
		"success_condition_value": "on",
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result)
	assert.Contains(t, result.Log, "Device is still off")

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_LambdaInvocationError(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock the lambda client to return an error
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return([]byte{}, assert.AnError)

	// Test config
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultPending, result.Result) // Result should remain pending on error
	assert.Contains(t, result.Log, "Error invoking lambda function")

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_InvalidJSONResponse(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock the lambda client to return invalid JSON
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return([]byte("invalid json"), nil)

	// Test config
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultPending, result.Result)

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_ComplexResponseData(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock complex response data similar to what you might get from your APIs
	responseData := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"devices": []map[string]interface{}{
				{
					"id":       "device-1",
					"name":     "HP Server",
					"status":   "on",
					"ip":       "192.168.1.10",
					"type":     "server",
					"location": "rack-1",
				},
				{
					"id":       "device-2",
					"name":     "Network Switch",
					"status":   "on",
					"ip":       "192.168.1.1",
					"type":     "switch",
					"location": "rack-1",
				},
			},
			"total_devices":  2,
			"online_devices": 2,
		},
		"timestamp": "2024-01-15T10:30:00Z",
		"message":   "Device discovery completed successfully",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config
	config := map[string]interface{}{
		"fn_name": "function/network-controller",
		"args": map[string]interface{}{
			"action": "discover_devices",
		},
		"success_condition_key":   "success",
		"success_condition_value": "true",
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Device discovery completed successfully")
	assert.Contains(t, result.Log, "HP Server")
	assert.Contains(t, result.Log, "Network Switch")

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_EmptySuccessCondition(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock response data
	responseData := map[string]interface{}{
		"status": "success",
		"data":   "test data",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config with empty success condition
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action": "test_action",
		},
		"success_condition_key":   "",
		"success_condition_value": "",
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result) // Should fail when success condition is empty

	mockLambdaClient.AssertExpectations(t)
}

func TestLambdaAgent_Execute_MissingSuccessCondition(t *testing.T) {
	// Setup
	mockLambdaClient := &MockLambdaFunction{}
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-lambda-step",
		},
		Result: entities.ResultPending,
	}

	agent := &LambdaAgent{
		step:         step,
		lambdaClient: mockLambdaClient,
	}

	// Mock response data
	responseData := map[string]interface{}{
		"status": "success",
		"data":   "test data",
	}
	responseBytes, _ := json.Marshal(responseData)

	// Mock the lambda client
	mockLambdaClient.On("Invoke", mock.AnythingOfType("client.LambdaFnConfig")).Return(responseBytes, nil)

	// Test config without success condition
	config := map[string]interface{}{
		"fn_name": "function/test-function",
		"args": map[string]interface{}{
			"action": "test_action",
		},
	}

	// Execute
	ctx := context.Background()
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result) // Should fail when success condition is missing

	mockLambdaClient.AssertExpectations(t)
}

// Test helper function to create a test step
func createTestStep(name string) *entities.Step {
	return &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: name,
		},
		Result: entities.ResultPending,
	}
}

// Test helper function to create test config
func createTestConfig(fnName string, args map[string]interface{}, successCondition map[string]string) map[string]interface{} {
	config := map[string]interface{}{
		"fn_name": fnName,
		"args":    args,
	}

	if successCondition != nil {
		config["success_condition_key"] = successCondition["key"]
		config["success_condition_value"] = successCondition["value"]
	}

	return config
}
