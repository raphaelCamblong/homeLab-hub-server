package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type PipelineHandler struct {
	pipelineUseCase usecase.PipelineUseCase
}

func NewPipelineHandler(pipelineUseCase usecase.PipelineUseCase) *PipelineHandler {
	return &PipelineHandler{pipelineUseCase: pipelineUseCase}
}

// ListPipelineTemplates godoc
// @Summary List all pipeline templates
// @Description Get all pipeline templates
// @Tags pipelines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.PipelineTemplate
// @Router /api/v1/pipelines [get]
func (h *PipelineHandler) ListPipelineTemplates(c *gin.Context) {
	templates, err := h.pipelineUseCase.ListPipelineTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetPipelineTemplate godoc
// @Summary Get a pipeline template by ID
// @Description Get a single pipeline template by its ID
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Template ID"
// @Security BearerAuth
// @Success 200 {object} entities.PipelineTemplate
// @Router /api/v1/pipelines/{id} [get]
func (h *PipelineHandler) GetPipelineTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.pipelineUseCase.GetPipelineTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// CreatePipelineTemplate godoc
// @Summary Create a new pipeline template
// @Description Create a new pipeline template
// @Tags pipelines
// @Accept json
// @Produce json
// @Param template body entities.PipelineTemplate true "Template object"
// @Security BearerAuth
// @Success 201 {object} entities.PipelineTemplate
// @Router /api/v1/pipelines [post]
func (h *PipelineHandler) CreatePipelineTemplate(c *gin.Context) {
	var template entities.PipelineTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.pipelineUseCase.CreatePipelineTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdatePipelineTemplate godoc
// @Summary Update a pipeline template
// @Description Update a pipeline template by ID
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Template ID"
// @Param template body entities.PipelineTemplate true "Template object"
// @Security BearerAuth
// @Success 200 {object} entities.PipelineTemplate
// @Router /api/v1/pipelines/{id} [put]
func (h *PipelineHandler) UpdatePipelineTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var template entities.PipelineTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	template.ID = uint(id)
	if err := h.pipelineUseCase.UpdatePipelineTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeletePipelineTemplate godoc
// @Summary Delete a pipeline template
// @Description Delete a pipeline template by ID
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Template ID"
// @Security BearerAuth
// @Success 200 {object} string "Template deleted successfully"
// @Router /api/v1/pipelines/{id} [delete]
func (h *PipelineHandler) DeletePipelineTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.pipelineUseCase.DeletePipelineTemplate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// CreateJob godoc
// @Summary Create a new job for a pipeline
// @Description Create a new execution job for a pipeline
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Pipeline ID"
// @Security BearerAuth
// @Success 201 {object} entities.Job
// @Router /api/v1/pipelines/{id}/jobs [post]
func (h *PipelineHandler) CreateJob(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	username := c.GetString("user_id")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	job, err := h.pipelineUseCase.CreateJob(uint(pipelineID), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// ListJobsByPipeline godoc
// @Summary List all jobs for a pipeline
// @Description Get all execution jobs for a specific pipeline
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Pipeline ID"
// @Security BearerAuth
// @Success 200 {array} entities.Job
// @Router /api/v1/pipelines/{id}/jobs [get]
func (h *PipelineHandler) ListJobsByPipeline(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	jobs, err := h.pipelineUseCase.ListJobsByPipeline(uint(pipelineID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// CompleteAllPipelineJobs godoc
// @Summary Complete all running jobs for a pipeline
// @Description Mark all running jobs as completed for a specific pipeline
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path uint true "Pipeline ID"
// @Security BearerAuth
// @Success 200 {object} string "Jobs completed successfully"
// @Router /api/v1/pipelines/{id}/jobs/complete [post]
func (h *PipelineHandler) CompleteAllPipelineJobs(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	if err := h.pipelineUseCase.CompleteAllPipelineJobs(uint(pipelineID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jobs completed successfully"})
}

// ListJobs godoc
// @Summary List all jobs across all pipelines
// @Description Get all execution jobs for all pipelines
// @Tags pipelines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Job
// @Router /api/v1/pipelines/jobs [get]
func (h *PipelineHandler) ListJobs(c *gin.Context) {
	jobs, err := h.pipelineUseCase.ListJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// StreamPipelinesJobs godoc
// @Summary Stream job updates
// @Description Get real-time updates for pipeline jobs via Server-Sent Events (SSE).
// The stream will automatically close when the job completes or when the client disconnects.
// @Tags pipelines
// @Accept json
// @Produce text/event-stream
// @Security BearerAuth
// @Success 200 {string} string "SSE stream of job updates"
// @Failure 500 {object} gin.H
// @Router /api/v1/pipelines/jobs/stream [get]
func (h *PipelineHandler) StreamPipelinesJobs(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	updates, err := h.pipelineUseCase.SubscribeToJobUpdates(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe to job updates"})
		return
	}

	logrus.Info("[Handler::StreamPipelinesJobs] Job stream opened")
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			logrus.Info("[Handler::StreamPipelinesJobs] Job stream closed")
			return false
		case job, ok := <-updates:
			if !ok {
				return false
			}
			c.SSEvent("job", job)
			logrus.Info("[Handler::StreamPipelinesJobs] Job stream updated")
			return true
		}
	})
}

// GetJob godoc
// @Summary Get a specific job
// @Description Get details of a specific job by its ID
// @Tags pipelines
// @Accept json
// @Produce json
// @Param jobId path uint true "Job ID"
// @Security BearerAuth
// @Success 200 {object} entities.Job
// @Router /api/v1/pipelines/jobs/{jobId} [get]
func (h *PipelineHandler) GetJob(c *gin.Context) {
	jobID, err := strconv.ParseUint(c.Param("jobId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.pipelineUseCase.GetJob(uint(jobID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// StopJob godoc
// @Summary Stop a specific job
// @Description Stop a specific job by its ID
// @Tags pipelines
// @Accept json
// @Produce json
// @Param jobId path uint true "Job ID"
// @Security BearerAuth
// @Success 200 {object} string "Job stopped successfully"
// @Router /api/v1/pipelines/jobs/{jobId}/stop [post]
func (h *PipelineHandler) StopJob(c *gin.Context) {
	jobID, err := strconv.ParseUint(c.Param("jobId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	if err := h.pipelineUseCase.StopJob(uint(jobID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job stopped successfully"})
}
