package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type KubeOperatorConfig struct {
	Namespace     string            `json:"namespace"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	Timeout       time.Duration     `json:"timeout"`
	Action        string            `json:"action"`
	KustomizePath string            `json:"path"`
}

type KubeAgent struct {
	step      *entities.Step
	k8sClient client.K8sClient
	cfg       *KubeOperatorConfig
}

func NewKubeAgent(step *entities.Step) Agent {
	k8sClient, err := client.NewK8sClient(&client.ClientOptions{})
	if err != nil {
		logrus.Fatalf("Failed to create k8s client: %v", err)
	}
	return &KubeAgent{step: step, k8sClient: *k8sClient}
}

func (e *KubeAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	logrus.Infof("Executing Kube Agent operator step: %s", e.step.StepTemplate.Name)

	err := e.loadConfig()
	if err != nil {
		return *e.step, err
	}

	switch e.cfg.Action {
	case "deploy":
		err = e.deploy()
	case "delete":
		err = e.delete()
	case "stop":
		err = e.stop()
	}

	if err != nil {
		e.step.Result = entities.ResultFailed
		e.step.Log += fmt.Sprintf("\n\t--|Error: %s", err.Error())
		return *e.step, err
	}

	return *e.step, nil
}

func (e *KubeAgent) loadConfig() error {
	cfg := KubeOperatorConfig{}
	if err := json.Unmarshal([]byte(e.step.StepTemplate.Config), &cfg); err != nil {
		return err
	}
	e.cfg = &cfg
	return nil
}

func (e *KubeAgent) deploy() error {
	result, err := e.k8sClient.DeployKustomize(e.cfg.KustomizePath, &client.KustomizeDeploymentOptions{
		Namespace:   e.cfg.Namespace,
		Labels:      e.cfg.Labels,
		Annotations: e.cfg.Annotations,
	})
	if err != nil {
		return err
	}
	e.addLog(result)

	if result.Success {
		return nil
	}
	return fmt.Errorf("deployment failed: %s", result.Message)
}

func (e *KubeAgent) stop() error {
	result, err := e.k8sClient.StopKustomize(e.cfg.KustomizePath, &client.KustomizeDeploymentOptions{
		Namespace:   e.cfg.Namespace,
		Labels:      e.cfg.Labels,
		Annotations: e.cfg.Annotations,
	})
	if err != nil {
		return err
	}

	e.addLog(result)

	if result.Success {
		return nil
	}
	return fmt.Errorf("stop failed: %s", result.Message)
}

func (e *KubeAgent) delete() error {
	result, err := e.k8sClient.DeleteKustomize(e.cfg.KustomizePath, &client.KustomizeDeploymentOptions{
		Namespace:   e.cfg.Namespace,
		Labels:      e.cfg.Labels,
		Annotations: e.cfg.Annotations,
	})
	if err != nil {
		return err
	}

	e.addLog(result)

	if result.Success {
		return nil
	}
	return fmt.Errorf("deletion failed: %s", result.Message)
}

func (e *KubeAgent) addLog(result *client.KustomizeDeploymentResult) {
	const resultTemplate = `
		--|Success: %t
		--|Message: %s
		--|DeployTime: %s
		`
	e.step.Log += fmt.Sprintf(resultTemplate, result.Success, result.Message, result.DeployTime)
}
