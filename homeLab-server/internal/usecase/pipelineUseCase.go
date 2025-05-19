package usecase

import (
	"context"

	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type PipelineUseCase interface {
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
	ListJobs() ([]entities.Job, error)
	ListJobsByPipeline(pipelineID uint) ([]entities.Job, error)
	GetJob(jobID uint) (*entities.Job, error)
	CreateJob(pipelineID uint, runBy string) (*entities.Job, error)
	CompleteJob(jobID uint, result string) error
	CompleteAllPipelineJobs(pipelineID uint) error
	StopJob(jobID uint) error

	// Step Management
	ListStepsByJob(jobID uint) ([]entities.Step, error)
	GetStep(stepID uint) (*entities.Step, error)
	UpdateStepStatus(stepID uint, status entities.Status, log string) error

	// Streaming
	SubscribeToJobUpdates(ctx context.Context) (<-chan entities.Job, error)
}

type pipelineUseCase struct {
	pipelineRepo repositories.PipelineRepository
}

func NewPipelineUseCase(pipelineRepo repositories.PipelineRepository) PipelineUseCase {
	return &pipelineUseCase{pipelineRepo: pipelineRepo}
}

// Pipeline Template Management implementations
func (u *pipelineUseCase) ListPipelineTemplates() ([]entities.PipelineTemplate, error) {
	return u.pipelineRepo.ListPipelineTemplates()
}

func (u *pipelineUseCase) GetPipelineTemplate(id uint) (*entities.PipelineTemplate, error) {
	return u.pipelineRepo.GetPipelineTemplate(id)
}

func (u *pipelineUseCase) CreatePipelineTemplate(template *entities.PipelineTemplate) error {
	return u.pipelineRepo.CreatePipelineTemplate(template)
}

func (u *pipelineUseCase) UpdatePipelineTemplate(template *entities.PipelineTemplate) error {
	return u.pipelineRepo.UpdatePipelineTemplate(template)
}

func (u *pipelineUseCase) DeletePipelineTemplate(id uint) error {
	return u.pipelineRepo.DeletePipelineTemplate(id)
}

// Step Template Management implementations
func (u *pipelineUseCase) ListStepTemplates() ([]entities.StepTemplate, error) {
	return u.pipelineRepo.ListStepTemplates()
}

func (u *pipelineUseCase) GetStepTemplate(id uint) (*entities.StepTemplate, error) {
	return u.pipelineRepo.GetStepTemplate(id)
}

func (u *pipelineUseCase) CreateStepTemplate(step *entities.StepTemplate) error {
	return u.pipelineRepo.CreateStepTemplate(step)
}

func (u *pipelineUseCase) UpdateStepTemplate(step *entities.StepTemplate) error {
	return u.pipelineRepo.UpdateStepTemplate(step)
}

func (u *pipelineUseCase) DeleteStepTemplate(id uint) error {
	return u.pipelineRepo.DeleteStepTemplate(id)
}

func (u *pipelineUseCase) AssociateStepsWithPipeline(pipelineID uint, stepIDs []uint) error {
	return u.pipelineRepo.AssociateStepsWithPipeline(pipelineID, stepIDs)
}

// Job Management implementations
func (u *pipelineUseCase) ListJobs() ([]entities.Job, error) {
	return u.pipelineRepo.ListJobs()
}

func (u *pipelineUseCase) ListJobsByPipeline(pipelineID uint) ([]entities.Job, error) {
	return u.pipelineRepo.ListJobsByPipeline(pipelineID)
}

func (u *pipelineUseCase) GetJob(jobID uint) (*entities.Job, error) {
	return u.pipelineRepo.GetJob(jobID)
}

func (u *pipelineUseCase) CreateJob(pipelineID uint, runBy string) (*entities.Job, error) {
	return u.pipelineRepo.CreateJob(pipelineID, runBy)
}

func (u *pipelineUseCase) CompleteJob(jobID uint, result string) error {
	return u.pipelineRepo.CompleteJob(jobID, result)
}

func (u *pipelineUseCase) CompleteAllPipelineJobs(pipelineID uint) error {
	jobs, err := u.pipelineRepo.ListJobsByPipeline(pipelineID)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if job.Status == entities.StatusRunning {
			if err := u.pipelineRepo.CompleteJob(job.ID, "Completed by system"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *pipelineUseCase) StopJob(jobID uint) error {
	return u.pipelineRepo.StopJob(jobID)
}

// Step Management implementations
func (u *pipelineUseCase) ListStepsByJob(jobID uint) ([]entities.Step, error) {
	return u.pipelineRepo.ListStepsByJob(jobID)
}

func (u *pipelineUseCase) GetStep(stepID uint) (*entities.Step, error) {
	return u.pipelineRepo.GetStep(stepID)
}

func (u *pipelineUseCase) UpdateStepStatus(stepID uint, status entities.Status, log string) error {
	return u.pipelineRepo.UpdateStepStatus(stepID, status, log)
}

// Streaming implementation
func (u *pipelineUseCase) SubscribeToJobUpdates(ctx context.Context) (<-chan entities.Job, error) {
	return u.pipelineRepo.SubscribeToJobUpdates(ctx)
}
