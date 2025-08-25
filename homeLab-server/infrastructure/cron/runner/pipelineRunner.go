package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron/agent"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

// JobUpdate represents an update for a job's status
type JobUpdate struct {
	Job   *entities.Job
	Error error
}

// StepUpdate represents an update for a step's status
type StepUpdate struct {
	Step      *entities.Step
	ParentJob *entities.Job
	Error     error
}

// PipelineRunner represents a single pipeline execution
type PipelineRunner struct {
	job          *entities.Job
	executors    map[string]agent.Agent
	cancelFunc   context.CancelFunc
	runningMutex sync.RWMutex
	jobChan      chan JobUpdate
	stepChan     chan StepUpdate
}

func NewPipelineRunner(job *entities.Job) *PipelineRunner {
	return &PipelineRunner{
		job:          job,
		executors:    make(map[string]agent.Agent),
		runningMutex: sync.RWMutex{},
		jobChan:      make(chan JobUpdate, 1),
		stepChan:     make(chan StepUpdate, len(job.Steps)),
	}
}

func (r *PipelineRunner) GetJobChannel() <-chan JobUpdate {
	return r.jobChan
}

func (r *PipelineRunner) GetStepChannel() <-chan StepUpdate {
	return r.stepChan
}

func (r *PipelineRunner) LoadExecutors() error {
	for i := range r.job.Steps {
		executor := agent.CreateAgent(&r.job.Steps[i])
		r.executors[r.job.Steps[i].StepTemplate.Name] = executor
	}
	return nil
}

func (r *PipelineRunner) RunJob(ctx context.Context) error {
	jobCtx, cancel := context.WithCancel(ctx)

	r.runningMutex.Lock()
	r.cancelFunc = cancel
	r.runningMutex.Unlock()

	go r.executeJob(jobCtx)
	return nil
}

func (r *PipelineRunner) executeJob(ctx context.Context) {
	defer r.handlePanic()
	defer close(r.jobChan)
	defer close(r.stepChan)

	r.setJobRunning()

	for i := range r.job.Steps {
		if err := r.executeStepWithContext(ctx, &r.job.Steps[i]); err != nil {
			r.handleStepError(&r.job.Steps[i], err)
			return
		}
	}

	r.handleJobCompletion()
}

func (r *PipelineRunner) executeStepWithContext(ctx context.Context, step *entities.Step) error {
	select {
	case <-ctx.Done():
		return r.handleCancellation(step)
	default:
		return r.ExecuteStep(step, ctx)
	}
}

func (r *PipelineRunner) ExecuteStep(step *entities.Step, ctx context.Context) error {
	config, err := r.parseStepConfig(step)
	if err != nil {
		return err
	}

	executor, ok := r.executors[step.StepTemplate.Name]
	if !ok {
		return fmt.Errorf("no executor found for step type: %s", step.StepTemplate.Name)
	}

	r.setStepRunning(step)
	return r.executeStepWithExecutor(executor, config, ctx)
}

func (r *PipelineRunner) parseStepConfig(step *entities.Step) (map[string]interface{}, error) {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(step.StepTemplate.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to parse step config: %v", err)
	}
	return config, nil
}

func (r *PipelineRunner) executeStepWithExecutor(
	executor agent.Agent,
	config map[string]interface{},
	ctx context.Context,
) error {
	step, err := executor.Execute(ctx, config)
	if err != nil {
		return fmt.Errorf("step %s failed: %v", step.StepTemplate.Name, err)
	}

	r.handleStepSuccess(&step)
	return nil
}

func (r *PipelineRunner) handleCancellation(step *entities.Step) error {
	now := time.Now()
	step.Status = entities.StatusStopped
	step.Result = entities.ResultSkipped
	step.EndedAt = &now
	r.stepChan <- StepUpdate{Step: step, ParentJob: r.job}

	r.job.Status = entities.StatusStopped
	r.job.Result = string(entities.ResultSkipped)
	r.job.EndedAt = &now
	r.jobChan <- JobUpdate{Job: r.job}

	return fmt.Errorf("job stopped")
}

func (r *PipelineRunner) handleStepError(step *entities.Step, err error) {
	step.Status = entities.StatusFailed
	step.Result = entities.ResultFailed
	r.stepChan <- StepUpdate{Step: step, ParentJob: r.job, Error: err}
}

func (r *PipelineRunner) setJobRunning() {
	now := time.Now()
	r.job.StartedAt = &now
	r.job.Status = entities.StatusRunning
	r.jobChan <- JobUpdate{Job: r.job}
}

func (r *PipelineRunner) setStepRunning(step *entities.Step) {
	step.Status = entities.StatusRunning
	now := time.Now()
	step.StartedAt = &now
	r.stepChan <- StepUpdate{Step: step, ParentJob: r.job}
}

func (r *PipelineRunner) handleStepSuccess(step *entities.Step) {
	step.Status = entities.StatusSuccess
	now := time.Now()
	step.EndedAt = &now
	r.stepChan <- StepUpdate{Step: step, ParentJob: r.job}
}

func (r *PipelineRunner) handleJobCompletion() {
	now := time.Now()
	r.job.EndedAt = &now
	r.job.Status = entities.StatusSuccess
	r.job.Result = fmt.Sprintf("Job completed successfully with %d steps", len(r.job.Steps))
	logrus.Infof("Job [%d] completed: %s", r.job.ID, r.job.Result)
	r.jobChan <- JobUpdate{Job: r.job}
}

func (r *PipelineRunner) handlePanic() {
	if rec := recover(); rec != nil {
		logrus.Errorf("Panic in job execution: %v", rec)
		now := time.Now()
		r.job.Status = entities.StatusFailed
		r.job.Result = fmt.Sprintf("Panic in job execution: %v", rec)
		r.job.EndedAt = &now
		r.jobChan <- JobUpdate{Job: r.job, Error: fmt.Errorf("panic: %v", rec)}
	}
}

func (r *PipelineRunner) StopJob() error {
	r.runningMutex.RLock()
	cancel := r.cancelFunc
	r.runningMutex.RUnlock()

	if cancel == nil {
		return fmt.Errorf("job %d is not running", r.job.ID)
	}

	cancel()
	return nil
}

func (r *PipelineRunner) GetJobStatus() *entities.Job {
	return r.job
}
