package cron

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron/runner"
)

type Cron struct {
	crons      map[string]*runner.PipelineRunner
	runningMux sync.RWMutex
}

func NewCron() *Cron {
	return &Cron{
		crons: make(map[string]*runner.PipelineRunner),
	}
}

func (c *Cron) AddCron(name string, cron *runner.PipelineRunner) {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()
	c.crons[name] = cron
}

func (c *Cron) GetCron(name string) (*runner.PipelineRunner, error) {
	c.runningMux.RLock()
	defer c.runningMux.RUnlock()

	cron, ok := c.crons[name]
	if !ok {
		return nil, fmt.Errorf("cron %s not found", name)
	}
	return cron, nil
}

func (c *Cron) StopCron(name string) error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	cron, ok := c.crons[name]
	if !ok {
		return fmt.Errorf("cron %s not found", name)
	}

	if err := cron.StopJob(); err != nil {
		return fmt.Errorf("failed to stop cron %s: %v", name, err)
	}

	delete(c.crons, name)
	return nil
}

func (c *Cron) ListRunningCrons() []string {
	c.runningMux.RLock()
	defer c.runningMux.RUnlock()

	running := make([]string, 0, len(c.crons))
	for name := range c.crons {
		running = append(running, name)
	}
	return running
}

func (c *Cron) StopAllCrons() error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	var lastErr error
	for name, cron := range c.crons {
		if err := cron.StopJob(); err != nil {
			logrus.Errorf("Failed to stop cron %s: %v", name, err)
			lastErr = err
		}
		delete(c.crons, name)
	}
	return lastErr
}

func (c *Cron) IsCronRunning(name string) bool {
	c.runningMux.RLock()
	defer c.runningMux.RUnlock()

	_, exists := c.crons[name]
	return exists
}
