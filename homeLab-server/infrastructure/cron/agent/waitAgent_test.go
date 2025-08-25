package agent

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

// MockAgent is a mock implementation of the Agent interface for testing
type MockAgent struct {
	mock.Mock
}

func (m *MockAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(entities.Step), args.Error(1)
}

func TestNewWaitAgent(t *testing.T) {
	// Test creating a new WaitAgent with an underlying agent
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
		},
	}
	mockAgent := &MockAgent{}

	agent := NewWaitAgent(step, mockAgent)
	assert.NotNil(t, agent)
	assert.IsType(t, &WaitAgent{}, agent)

	waitAgent := agent.(*WaitAgent)
	assert.Equal(t, step, waitAgent.step)
	assert.Equal(t, mockAgent, waitAgent.agent)
}

func TestNewWaitAgentEmpty(t *testing.T) {
	// Test creating a new WaitAgent without an underlying agent
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
		},
	}

	agent := NewWaitAgentEmpty(step)
	assert.NotNil(t, agent)
	assert.IsType(t, &WaitAgent{}, agent)

	waitAgent := agent.(*WaitAgent)
	assert.Equal(t, step, waitAgent.step)
	assert.Nil(t, waitAgent.agent)
}

func TestWaitAgent_Execute_SuccessOnFirstTry(t *testing.T) {
	// Setup
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 5,
				"retry": 3
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock successful execution on first try
	successStep := entities.Step{
		Result: entities.ResultSuccess,
		Log:    "Success on first try",
	}
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(successStep, nil)

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusSuccess, result.Status)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Trying to execute agent")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_SuccessAfterRetries(t *testing.T) {
	// Setup
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 1,
				"retry": 3
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock failed executions first, then success
	failedStep := entities.Step{
		Result: entities.ResultFailed,
		Log:    "Failed attempt",
	}
	successStep := entities.Step{
		Result: entities.ResultSuccess,
		Log:    "Success after retries",
	}

	// First two calls fail, third succeeds
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(failedStep, nil).Twice()
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(successStep, nil).Once()

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusSuccess, result.Status)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Trying to execute agent")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_FailureAfterMaxRetries(t *testing.T) {
	// Setup
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 1,
				"retry": 2
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock all executions to fail
	failedStep := entities.Step{
		Result: entities.ResultFailed,
		Log:    "Failed attempt",
	}
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(failedStep, nil).Times(3) // 2 retries + 1 initial

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result)
	assert.Contains(t, err.Error(), "wait step failed after 2 attempts")
	assert.Contains(t, result.Log, "Wait step failed after 2 attempts")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_Timeout(t *testing.T) {
	// Setup - start time in the past to simulate timeout
	pastTime := time.Now().Add(-2 * time.Minute)
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 5,
				"retry": 10
			}`,
		},
		StartedAt: &pastTime,
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock failed execution
	failedStep := entities.Step{
		Result: entities.ResultFailed,
		Log:    "Failed attempt",
	}
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(failedStep, nil).Once()

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result)
	assert.Contains(t, err.Error(), "wait step timed out after 60 seconds")
	assert.Contains(t, result.Log, "Wait step timed out after 60 seconds")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_AgentError(t *testing.T) {
	// Setup
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 1,
				"retry": 2
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock agent to return an error
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(entities.Step{}, assert.AnError).Times(3)

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result)
	assert.Contains(t, err.Error(), "wait step failed after 2 attempts")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_NoAgent(t *testing.T) {
	// Setup - WaitAgent without underlying agent
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 60,
				"interval": 1,
				"retry": 2
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	agent := &WaitAgent{
		step:  step,
		agent: nil, // No underlying agent
	}

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultFailed, result.Result)
	assert.Contains(t, err.Error(), "wait step failed after 2 attempts")
}

func TestWaitAgent_Execute_DefaultConfig(t *testing.T) {
	// Setup - step without config to test defaults
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name:   "test-wait-step",
			Config: `{}`, // Empty config
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock successful execution
	successStep := entities.Step{
		Result: entities.ResultSuccess,
		Log:    "Success",
	}
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(successStep, nil).Once()

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusSuccess, result.Status)
	assert.Equal(t, entities.ResultSuccess, result.Result)

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_Execute_InvalidConfig(t *testing.T) {
	// Setup - step with invalid JSON config
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name:   "test-wait-step",
			Config: `{ invalid json }`, // Invalid JSON
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Execute
	ctx := context.Background()
	config := map[string]interface{}{
		"action": "test_action",
	}
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, entities.ResultPending, result.Result)
}

func TestWaitAgent_RealWorldExample(t *testing.T) {
	// Test with a real-world example from your pipeline config
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "[Network-Controller + W] Turn On Server Device plug",
			Config: `{
				"timeout": 300,
				"interval": 5,
				"retry": 3
			}`,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}

	mockAgent := &MockAgent{}
	agent := &WaitAgent{
		step:  step,
		agent: mockAgent,
	}

	// Mock the underlying lambda agent behavior
	// First two attempts fail (device still off), third succeeds (device on)
	failedStep := entities.Step{
		Result: entities.ResultFailed,
		Log:    "Device status: off",
	}
	successStep := entities.Step{
		Result: entities.ResultSuccess,
		Log:    "Device status: on",
	}

	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(failedStep, nil).Twice()
	mockAgent.On("Execute", mock.Anything, mock.Anything).Return(successStep, nil).Once()

	// Execute with config similar to your pipeline
	ctx := context.Background()
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
	result, err := agent.Execute(ctx, config)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusSuccess, result.Status)
	assert.Equal(t, entities.ResultSuccess, result.Result)
	assert.Contains(t, result.Log, "Trying to execute agent")

	mockAgent.AssertExpectations(t)
}

func TestWaitAgent_LoadConfig(t *testing.T) {
	// Test the loadConfig method directly
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 120,
				"interval": 10,
				"retry": 5
			}`,
		},
	}

	agent := &WaitAgent{
		step: step,
	}

	err := agent.loadConfig()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 120, agent.cfg.Timeout)
	assert.Equal(t, 10, agent.cfg.Interval)
	assert.Equal(t, 5, agent.cfg.Retry)
	assert.Equal(t, 120*time.Second, agent.cfg.timeout)
	assert.Equal(t, 10*time.Second, agent.cfg.interval)
}

func TestWaitAgent_LoadConfig_Defaults(t *testing.T) {
	// Test that defaults are applied when values are missing or zero
	step := &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name: "test-wait-step",
			Config: `{
				"timeout": 0,
				"interval": 0,
				"retry": 0
			}`,
		},
	}

	agent := &WaitAgent{
		step: step,
	}

	err := agent.loadConfig()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int(7*time.Minute/time.Second), agent.cfg.Timeout) // Default timeout
	assert.Equal(t, 5, agent.cfg.Interval)                             // Default interval
	assert.Equal(t, 4, agent.cfg.Retry)                                // Default retry
}

// Test helper function to create a test wait step
func createTestWaitStep(name, config string) *entities.Step {
	return &entities.Step{
		StepTemplate: &entities.StepTemplate{
			Name:   name,
			Config: config,
		},
		StartedAt: &[]time.Time{time.Now()}[0],
		Result:    entities.ResultPending,
	}
}

// Test helper function to create wait config
func createWaitConfig(timeout, interval, retry int) string {
	config := map[string]interface{}{
		"timeout":  timeout,
		"interval": interval,
		"retry":    retry,
	}

	configBytes, _ := json.Marshal(config)
	return string(configBytes)
}
