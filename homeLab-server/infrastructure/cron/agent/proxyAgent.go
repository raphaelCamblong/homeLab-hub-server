package agent

import (
	"context"
	"fmt"
	"time"

	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

// --- AgentProxy ---
// AgentProxy is a proxy for the Agent interface.
// Used to add logging, timing and status to the agent execution.
type AgentProxy struct {
	step  *entities.Step
	agent Agent
}

func NewAgentProxy(step *entities.Step, agent Agent) Agent {
	return &AgentProxy{step: step, agent: agent}
}

func (e *AgentProxy) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	now := time.Now()
	e.step.StartedAt = &now
	e.step.Status = entities.StatusRunning
	e.step.Log += "\nStep started"
	_, err := e.agent.Execute(ctx, config)
	endnow := time.Now()
	e.step.EndedAt = &endnow
	if err != nil {
		e.step.Status = entities.StatusFailed
		e.step.Log += fmt.Sprintf("\nError: %v", err)
		e.step.Result = entities.ResultFailed
		return *e.step, err
	}
	e.step.Status = entities.StatusCompleted
	e.step.Result = entities.ResultSuccess
	return *e.step, nil
}

// --- WaitWrapperAgent ---
type WaitWrapperAgent struct {
	step      *entities.Step
	wait      Agent
	validator Agent
}

func NewWaitWrapperAgent(step *entities.Step, agent Agent) Agent {
	return &WaitWrapperAgent{step: step}
}

func (e *WaitWrapperAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	e.setUpAgent()
	return e.wait.Execute(ctx, config)
}

func (e *WaitWrapperAgent) setUpAgent() error {
	e.validator = NewLambdaAgent(e.step)
	e.wait = NewWaitAgent(e.step, e.validator)
	return nil
}
