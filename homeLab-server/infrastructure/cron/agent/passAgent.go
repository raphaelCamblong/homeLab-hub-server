package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type PassAgent struct {
	step *entities.Step
}

func NewPassAgent(step *entities.Step) Agent {
	return &PassAgent{step: step}
}

func (e *PassAgent) Execute(ctx context.Context, config map[string]interface{}) (entities.Step, error) {
	logrus.Infof("Test step %s running...", e.step.StepTemplate.Name)
	time.Sleep(5 * time.Second)
	e.step.Result = "Test step completed"
	e.step.Log += "\nTest step completed"
	return *e.step, nil
}
