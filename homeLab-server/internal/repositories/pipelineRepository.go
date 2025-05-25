package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron/runner"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/streaming"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type PipelineRepository interface {
	// Pipeline Template Management
	ListPipelineTemplates() ([]entities.PipelineTemplate, error)
	GetPipelineTemplate(id uint) (*entities.PipelineTemplate, error)
	CreatePipelineTemplate(template *entities.PipelineTemplate) error
	UpdatePipelineTemplate(template *entities.PipelineTemplate) error
	DeletePipelineTemplate(id uint) error

	// Step Template Management
	ListStepTemplates() ([]entities.StepTemplate, error)
	GetStepTemplate(id uint) (*entities.StepTemplate, error)
	CreateStepTemplate(step *entities.StepTemplate) error
	UpdateStepTemplate(step *entities.StepTemplate) error
	DeleteStepTemplate(id uint) error
	AssociateStepsWithPipeline(pipelineID uint, stepIDs []uint) error

	// Job Management
	ListJobs(status *entities.Status) ([]entities.Job, error)
	ListJobsByPipeline(pipelineID uint) ([]entities.Job, error)
	GetJob(jobID uint) (*entities.Job, error)
	CreateJob(pipelineID uint, runBy string) (*entities.Job, error)
	UpdateJobStatus(jobID uint, status entities.Status) error
	UpdateJob(job *entities.Job) error
	CompleteJob(jobID uint, result string) error
	StopJob(jobID uint) error
	DeleteJob(job *entities.Job) error

	// Step Management
	ListStepsByJob(jobID uint) ([]entities.Step, error)
	GetStep(stepID uint) (*entities.Step, error)
	UpdateStepStatus(stepID uint, status entities.Status, log string) error
	UpdateStep(step *entities.Step) error

	// Streaming
	SubscribeToJobUpdates(ctx context.Context) (<-chan entities.Job, error)
}

type pipelineRepository struct {
	db          *gorm.DB
	streamHub   *streaming.StreamHub
	cronManager *cron.Cron
}

func NewPipelineRepository(db database.Database, streamHub *streaming.StreamHub, cronManager *cron.Cron) PipelineRepository {
	return &pipelineRepository{
		db:          db.Get(),
		streamHub:   streamHub,
		cronManager: cronManager,
	}
}

// Pipeline Template Management implementations
func (r *pipelineRepository) ListPipelineTemplates() ([]entities.PipelineTemplate, error) {
	var templates []entities.PipelineTemplate
	result := r.db.Preload("Steps").
		Preload("RunningJobs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Find(&templates)
	return templates, result.Error
}

func (r *pipelineRepository) GetPipelineTemplate(id uint) (*entities.PipelineTemplate, error) {
	var template entities.PipelineTemplate
	result := r.db.Preload("Steps").
		Preload("RunningJobs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&template, id)
	return &template, result.Error
}

func (r *pipelineRepository) CreatePipelineTemplate(template *entities.PipelineTemplate) error {
	if err := r.db.Create(template).Error; err != nil {
		return err
	}
	return nil
}

func (r *pipelineRepository) UpdatePipelineTemplate(template *entities.PipelineTemplate) error {
	if err := r.db.Save(template).Error; err != nil {
		return err
	}
	return nil
}

func (r *pipelineRepository) DeletePipelineTemplate(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pipeline entities.PipelineTemplate
		if err := tx.Preload("Steps").First(&pipeline, id).Error; err != nil {
			return err
		}

		return tx.Select("Steps").Delete(&pipeline).Error
	})
}

// Step Template Management implementations
func (r *pipelineRepository) ListStepTemplates() ([]entities.StepTemplate, error) {
	var templates []entities.StepTemplate
	return templates, r.db.Find(&templates).Error
}

func (r *pipelineRepository) GetStepTemplate(id uint) (*entities.StepTemplate, error) {
	var template entities.StepTemplate
	return &template, r.db.First(&template, id).Error
}

func (r *pipelineRepository) CreateStepTemplate(step *entities.StepTemplate) error {
	return r.db.Create(step).Error
}

func (r *pipelineRepository) UpdateStepTemplate(step *entities.StepTemplate) error {
	return r.db.Save(step).Error
}

func (r *pipelineRepository) UpdateStep(step *entities.Step) error {
	return r.db.Save(step).Error
}

func (r *pipelineRepository) DeleteStepTemplate(id uint) error {
	return r.db.Delete(&entities.StepTemplate{}, id).Error
}

func (r *pipelineRepository) AssociateStepsWithPipeline(pipelineID uint, stepIDs []uint) error {
	var steps []entities.StepTemplate
	if err := r.db.Find(&steps, stepIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entities.PipelineTemplate{Model: gorm.Model{ID: pipelineID}}).
		Association("Steps").
		Replace(steps)
}

// Job Management implementations
func (r *pipelineRepository) ListJobs(status *entities.Status) ([]entities.Job, error) {
	var jobs []entities.Job
	query := r.db.Preload("Steps").Preload("Pipeline").Preload("Steps.StepTemplate")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	return jobs, query.Find(&jobs).Error
}

func (r *pipelineRepository) ListJobsByPipeline(pipelineID uint) ([]entities.Job, error) {
	var jobs []entities.Job
	return jobs, r.db.Preload("Steps").
		Where("pipeline_id = ?", pipelineID).
		Order("created_at desc").
		Find(&jobs).Error
}

func (r *pipelineRepository) GetJob(jobID uint) (*entities.Job, error) {
	var job entities.Job
	return &job, r.db.Preload("Steps").Preload("Pipeline").Preload("Steps.StepTemplate").First(&job, jobID).Error
}

func (r *pipelineRepository) CreateJob(pipelineID uint, runBy string) (*entities.Job, error) {
	var job *entities.Job
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var pipeline entities.PipelineTemplate
		if err := tx.Preload("Steps").First(&pipeline, pipelineID).Error; err != nil {
			return err
		}

		now := time.Now()
		job = &entities.Job{
			PipelineID: pipelineID,
			Status:     entities.StatusPending,
			RunBy:      runBy,
			StartedAt:  &now,
		}

		if err := tx.Create(job).Error; err != nil {
			return err
		}

		steps := make([]entities.Step, len(pipeline.Steps))
		for i, stepTemplate := range pipeline.Steps {
			steps[i] = entities.Step{
				JobID:          job.ID,
				StepTemplateID: stepTemplate.ID,
				Status:         entities.StatusPending,
				ExecutionOrder: i + 1,
			}
		}

		if err := tx.Create(&steps).Error; err != nil {
			return err
		}

		if err := tx.Preload("Steps.StepTemplate").Preload("Pipeline").First(job, job.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	var err error
	if err = r.scheduleJob(job); err != nil {
		logrus.Errorf("Failed to start job %d: %v", job.ID, err)
		job.Status = entities.StatusFailed
		job.Result = err.Error()
		now := time.Now()
		job.EndedAt = &now
		err = r.UpdateJob(job)
	}
	return job, err
}

func (r *pipelineRepository) scheduleJob(job *entities.Job) error {
	runner := runner.NewPipelineRunner(job)
	r.cronManager.AddCron(fmt.Sprintf("job-%d", job.ID), runner)

	if err := runner.LoadExecutors(); err != nil {
		return err
	}

	go r.handleJobUpdates(runner)
	go r.handleStepUpdates(runner.GetStepChannel())

	if err := runner.RunJob(context.Background()); err != nil {
		return err
	}

	return nil
}

func (r *pipelineRepository) handleJobCompletion(job *entities.Job) {
	r.streamHub.Publish(streaming.Event{
		Type: streaming.JobEvent,
		Data: *job,
	})

	r.StopJob(job.ID)
	// r.DeleteJob(job)
}

func (r *pipelineRepository) handleJobUpdates(runner *runner.PipelineRunner) {
	jobChan := runner.GetJobChannel()
	for update := range jobChan {
		if update.Error != nil {
			logrus.Errorf("Job %d error: %v", update.Job.ID, update.Error)
		}
		if err := r.UpdateJob(update.Job); err != nil {
			logrus.Errorf("Failed to save job update: %v", err)
		}

		if update.Job.IsCompleted() {
			r.handleJobCompletion(update.Job)
		} else {
			r.streamHub.Publish(streaming.Event{
				Type: streaming.JobEvent,
				Data: *update.Job,
			})
		}
	}
}

func (r *pipelineRepository) handleStepUpdates(stepChan <-chan runner.StepUpdate) {
	for update := range stepChan {
		if update.Error != nil {
			logrus.Errorf("Step %d error: %v", update.Step.ID, update.Error)
		}
		if err := r.UpdateStep(update.Step); err != nil {
			logrus.Errorf("Failed to save step update: %v", err)
		}
		r.streamHub.Publish(streaming.Event{
			Type: streaming.JobEvent,
			Data: *update.ParentJob,
		})
	}
}

func (r *pipelineRepository) UpdateJobStatus(jobID uint, status entities.Status) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
	}

	if status == entities.StatusCompleted || status == entities.StatusFailed {
		updates["ended_at"] = &now
	}

	return r.db.Model(&entities.Job{}).Where("id = ?", jobID).Updates(updates).Error
}

func (r *pipelineRepository) UpdateJob(job *entities.Job) error {
	return r.db.Save(job).Error
}

func (r *pipelineRepository) StopJob(jobID uint) error {
	return r.cronManager.StopCron(fmt.Sprintf("job-%d", jobID))
}

func (r *pipelineRepository) DeleteJob(job *entities.Job) error {
	result := r.db.Delete(&job) //.Unscoped()
	if result.Error != nil {
		return fmt.Errorf("failed to delete job: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("job %d not found", job.ID)
	}
	return nil
}

func (r *pipelineRepository) CompleteJob(jobID uint, result string) error {
	now := time.Now()
	return r.db.Model(&entities.Job{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":   entities.StatusCompleted,
		"result":   result,
		"ended_at": &now,
	}).Error
}

// Step Management implementations
func (r *pipelineRepository) ListStepsByJob(jobID uint) ([]entities.Step, error) {
	var steps []entities.Step
	return steps, r.db.Preload("StepTemplate").
		Where("job_id = ?", jobID).
		Order("execution_order asc").
		Find(&steps).Error
}

func (r *pipelineRepository) GetStep(stepID uint) (*entities.Step, error) {
	var step entities.Step
	logrus.Info("Fetching step with ID:", stepID)
	return &step, r.db.Preload("StepTemplate").First(&step, stepID).Error
}

func (r *pipelineRepository) UpdateStepStatus(stepID uint, status entities.Status, log string) error {
	var step entities.Step
	if err := r.db.Preload("Job").First(&step, stepID).Error; err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
		"log":    log,
	}

	if status == entities.StatusRunning {
		updates["started_at"] = &now
	}

	if status == entities.StatusSuccess || status == entities.StatusFailed {
		updates["ended_at"] = &now
	}

	if err := r.db.Model(&step).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// Streaming implementation

func (r *pipelineRepository) SubscribeToJobUpdates(ctx context.Context) (<-chan entities.Job, error) {
	eventCh, err := r.streamHub.Subscribe(ctx, streaming.JobEvent)
	if err != nil {
		return nil, err
	}

	jobCh := make(chan entities.Job)

	go func() {
		defer close(jobCh)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventCh:
				if !ok {
					return
				}
				if job, ok := event.Data.(entities.Job); ok {
					jobCh <- job
				}
			}
		}
	}()

	return jobCh, nil
}
