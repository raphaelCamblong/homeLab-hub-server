package agent

import (
	"context"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type Agent interface {
	Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error)
}

func CreateAgent(step *entities.Step) Agent {
	agents := map[string]Agent{
		"lambda+wait": NewWaitWrapperAgent(step, NewLambdaAgent(step)),
		"lambda":      NewLambdaAgent(step),
		"wait":        NewWaitAgentEmpty(step),
		"pass":        NewPassAgent(step),
	}

	agent, exists := agents[step.StepTemplate.Type]
	if !exists {
		logrus.Warnf("No agent found for step type: %s, using default pass agent", step.StepTemplate.Name)
		agent = agents["pass"]
	}

	return NewAgentProxy(step, agent)
}
