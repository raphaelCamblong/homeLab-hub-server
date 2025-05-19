package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type StepExecutor interface {
	Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error)
}

func CreateExecutor(step *entities.Step) StepExecutor {
	executors := map[string]StepExecutor{
		// "lambda": NewLambdaStepExecutor(step),
		"lambda": NewTestStepExecutor(step),
	}

	executor, exists := executors[step.StepTemplate.Type]
	if !exists {
		logrus.Warnf("No executor found for step type: %s, using default lambda executor", step.StepTemplate.Name)
		executor = executors["lambda"]
	}

	return NewStepExecutorProxy(step, executor)
}

func NewStepExecutorProxy(step *entities.Step, executor StepExecutor) StepExecutor {
	return &StepExecutorProxy{step: step, executor: executor}
}

type StepExecutorProxy struct {
	step     *entities.Step
	executor StepExecutor
}

func (e *StepExecutorProxy) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	now := time.Now()
	e.step.StartedAt = &now
	e.step.Status = entities.StatusRunning
	e.step.Log = "Step started"
	_, err := e.executor.Execute(ctx, config)
	endnow := time.Now()
	e.step.EndedAt = &endnow
	if err != nil {
		e.step.Status = entities.StatusFailed
		e.step.Result = err.Error()
		return *e.step, err
	}
	e.step.Status = entities.StatusSuccess
	// e.step.Result = step.Result
	// e.step.Log = step.Log
	return *e.step, nil
}

// --LambdaStepExecutor--
type LambdaStepExecutor struct {
	step         *entities.Step
	lambdaClient client.LambdaFunction
}

func (e *LambdaStepExecutor) checkConfig(config map[string]interface{}) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if config["args"] == nil {
		return fmt.Errorf("args is nil")
	}
	if config["fn_name"] == nil {
		return fmt.Errorf("fn_name is nil")
	}
	return nil
}

func (e *LambdaStepExecutor) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	logrus.Infof("Executing lambda step: %s", e.step.StepTemplate.Name)

	if err := e.checkConfig(config); err != nil {
		return *e.step, err
	}
	payload, err := json.Marshal(config["args"])
	if err != nil {
		return *e.step, err
	}
	resp, err := e.lambdaClient.Invoke(config["fn_name"].(string), payload)
	if err != nil {
		e.step.Status = entities.StatusFailed
		e.step.Result = err.Error()
		return *e.step, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp, &response); err != nil {
		e.step.Status = entities.StatusFailed
		e.step.Result = err.Error()
		return *e.step, err
	}
	e.step.Result = fmt.Sprintf("%v", response)
	e.step.Log = fmt.Sprintf("Lambda step %s completed with response: %v", e.step.StepTemplate.Name, response)
	return *e.step, nil
}

func NewLambdaStepExecutor(step *entities.Step) StepExecutor {
	cfg := config.Get()
	lambdaClient := client.NewLambdaFunction(cfg.Client.Lambda, client.NewHttpClient(20*time.Second))
	return &LambdaStepExecutor{step: step, lambdaClient: lambdaClient}
}

// --Test StepExecutor--
type TestStepExecutor struct {
	step *entities.Step
}

func (e *TestStepExecutor) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	now := time.Now()
	e.step.StartedAt = &now
	logrus.Infof("Test step %s running...", e.step.StepTemplate.Name)
	time.Sleep(5 * time.Second)
	e.step.Status = entities.StatusSuccess
	e.step.Result = "Test step completed"
	e.step.Log = "Test step completed"
	e.step.EndedAt = &now
	return *e.step, nil
}

func NewTestStepExecutor(step *entities.Step) StepExecutor {
	return &TestStepExecutor{step: step}
}
