package routes

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func PipelinesRoutes(infra *infrastructure.Infrastructure, pipelineRepo repositories.PipelineRepository) {
	pipelineUseCase := usecase.NewPipelineUseCase(pipelineRepo)
	pipelineHandler := handlers.NewPipelineHandler(pipelineUseCase)

	v1 := infra.Router.Get().Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		pipelines := v1.Group("/pipelines")
		{
			// Pipeline Template endpoints
			pipelines.GET("", middleware.HasPermission("pipelines", "read"), pipelineHandler.ListPipelineTemplates)
			pipelines.GET("/:id", middleware.HasPermission("pipelines", "read"), pipelineHandler.GetPipelineTemplate)
			pipelines.POST("", middleware.HasPermission("admin", "create"), pipelineHandler.CreatePipelineTemplate)
			pipelines.PUT("/:id", middleware.HasPermission("admin", "update"), pipelineHandler.UpdatePipelineTemplate)
			pipelines.DELETE("/:id", middleware.HasPermission("admin", "delete"), pipelineHandler.DeletePipelineTemplate)

			// Global Job Management endpoints
			jobs := pipelines.Group("/jobs")
			{
				jobs.GET("", middleware.HasPermission("pipelines", "read"), pipelineHandler.ListJobs)
				jobs.GET("/stream", middleware.HasPermission("pipelines", "read"), pipelineHandler.StreamPipelinesJobs)
				jobs.GET("/:jobId", middleware.HasPermission("pipelines", "read"), pipelineHandler.GetJob)
				jobs.POST("/:jobId/stop", middleware.HasPermission("pipelines", "delete"), pipelineHandler.StopJob)
			}

			// Pipeline-specific Job endpoints
			pipelineJobs := pipelines.Group("/:id/jobs")
			{
				pipelineJobs.GET("", middleware.HasPermission("pipelines", "read"), pipelineHandler.ListJobsByPipeline)
				pipelineJobs.POST("", middleware.HasPermission("pipelines", "execute"), pipelineHandler.CreateJob)
				pipelineJobs.POST("/complete", middleware.HasPermission("pipelines", "execute"), pipelineHandler.CompleteAllPipelineJobs)
			}
		}
	}
}
