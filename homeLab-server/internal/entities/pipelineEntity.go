package entities

import (
	"time"

	"gorm.io/gorm"
)

type PipelineTemplate struct {
	gorm.Model

	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`

	// Relationships
	Steps       []StepTemplate `json:"steps" gorm:"many2many:pipeline_step_templates;"`
	RunningJobs []Job          `json:"running_jobs" gorm:"foreignKey:PipelineID;constraint:OnDelete:CASCADE"`
}

type StepTemplate struct {
	gorm.Model

	Type        string `json:"type" gorm:"type:varchar(255);not null"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	Config      string `json:"config" gorm:"type:text"`
}

type Job struct {
	gorm.Model

	PipelineID uint       `json:"pipeline_id" gorm:"not null"`
	Status     Status     `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	RunBy      string     `json:"run_by" gorm:"type:varchar(255)"`
	StartedAt  *time.Time `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"`
	Result     string     `json:"result" gorm:"type:text"`

	// Relationships
	Pipeline PipelineTemplate `json:"pipeline" gorm:"foreignKey:PipelineID"`
	Steps    []Step           `json:"steps" gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}

type Step struct {
	gorm.Model

	JobID          uint       `json:"job_id" gorm:"not null;index"`
	StepTemplateID uint       `json:"step_template_id" gorm:"not null;index"`
	Status         Status     `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Log            string     `json:"log" gorm:"type:text"`
	Result         string     `json:"result" gorm:"type:text"`
	StartedAt      *time.Time `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at"`
	ExecutionOrder int        `json:"execution_order" gorm:"not null"`

	// Relationships
	Job          *Job          `json:"-" gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	StepTemplate *StepTemplate `json:"step_template" gorm:"foreignKey:StepTemplateID;references:ID"`
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusStopped   Status = "stopped"
)

// Helper functions
func (j *Job) IsCompleted() bool {
	return j.Status == StatusCompleted || j.Status == StatusSuccess
}

func (j *Job) IsFailed() bool {
	return j.Status == StatusFailed
}

func (j *Job) IsStopped() bool {
	return j.Status == StatusStopped
}

func (j *Job) IsRunning() bool {
	return j.Status == StatusRunning
}

func (j *Job) IsPending() bool {
	return j.Status == StatusPending
}

func (j *Job) IsSuccess() bool {
	return j.Status == StatusSuccess
}
