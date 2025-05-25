package agent

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

type LambdaAgent struct {
	step         *entities.Step
	lambdaClient client.LambdaFunction
}

func NewLambdaAgent(step *entities.Step) Agent {
	cfg := config.Get()
	lambdaClient := client.NewLambdaFunction(cfg.Client.Lambda, client.NewHttpClient(20*time.Second))
	return &LambdaAgent{step: step, lambdaClient: lambdaClient}
}

func (e *LambdaAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	logrus.Infof("Executing lambda step: %s", e.step.StepTemplate.Name)

	lambdaFnCfg, err := client.NewLambdaFnConfig(config)
	if err != nil {
		return *e.step, err
	}
	e.step.Log += fmt.Sprintf("\nInvoking lambda function: %s", e.step.StepTemplate.Name)
	resp, err := e.lambdaClient.Invoke(*lambdaFnCfg)
	if err != nil {
		e.step.Log += fmt.Sprintf("\nError invoking lambda function: %v", err)
		return *e.step, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp, &response); err != nil {
		return *e.step, err
	}
	e.step.Log += fmt.Sprintf("\n%v", response)
	e.step.Result = fmt.Sprintf("Lambda step %s completed with response", e.step.StepTemplate.Name)
	return *e.step, nil
}
