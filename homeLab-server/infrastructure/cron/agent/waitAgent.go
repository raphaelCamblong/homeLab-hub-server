package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type ConfigWaitStep struct {
	Timeout   int    `json:"timeout"`
	Interval  int    `json:"interval"`
	Count     int    `json:"count"`
	Condition string `json:"condition"`
	timeout   time.Duration
	interval  time.Duration
}

type WaitAgent struct {
	cfg   ConfigWaitStep
	step  *entities.Step
	agent Agent
}

func NewWaitAgent(step *entities.Step, agent Agent) Agent {
	return &WaitAgent{step: step, agent: agent}
}

func NewWaitAgentEmpty(step *entities.Step) Agent {
	return &WaitAgent{step: step}
}

func (e *WaitAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	err := e.loadConfig()
	if err != nil {
		return *e.step, err
	}

	for {
		if e.checkCount() {
			e.step.Log += "\nWait step failed after %d attempts"
			return *e.step, fmt.Errorf("wait step failed after %d attempts", e.cfg.Count)
		}
		if e.checkTimeout() {
			e.step.Log += "\nWait step timed out after %d seconds"
			return *e.step, fmt.Errorf("wait step timed out after %d seconds", e.cfg.Timeout)
		}
		if e.try(ctx, config) {
			break
		}
		e.sleep()
		e.cfg.Count--
	}
	e.step.Status = entities.StatusSuccess
	e.step.Result = "Wait step completed successfully"
	return *e.step, nil
}

func (e *WaitAgent) loadConfig() error {
	if err := json.Unmarshal([]byte(e.step.StepTemplate.Config), &e.cfg); err != nil {
		return err
	}
	if e.cfg.Count <= 0 {
		e.cfg.Count = 4
	}
	if e.cfg.Interval <= 0 {
		e.cfg.Interval = 5
	}
	if e.cfg.Timeout <= 0 {
		e.cfg.Timeout = int(5 * time.Minute / time.Second)
	}
	e.cfg.timeout = time.Duration(e.cfg.Timeout) * time.Second
	e.cfg.interval = time.Duration(e.cfg.Interval) * time.Second
	return nil
}

func (e *WaitAgent) try(ctx context.Context, config map[string]interface{}) bool {
	if e.agent == nil {
		return false
	}
	e.step.Log += "\nTrying to execute agent"
	_, err := e.agent.Execute(ctx, config)
	return err == nil
}

func (e *WaitAgent) sleep() error {
	time.Sleep(e.cfg.interval)
	return nil
}

func (e *WaitAgent) checkCount() bool {
	return e.cfg.Count > 0
}

func (e *WaitAgent) checkTimeout() bool {
	return time.Since(*e.step.StartedAt) > e.cfg.timeout
}
