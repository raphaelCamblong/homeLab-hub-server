package routes

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func ServiceRoutes(infra *infrastructure.Infrastructure, serviceRepo repositories.ServiceRepository) {
	serviceUseCase := usecase.NewServiceUseCase(serviceRepo)
	serviceHandler := handlers.NewServiceHandler(serviceUseCase)

	protected := infra.Router.Get().Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		services := protected.Group("/services")
		{
			// List and get services - requires services:read permission
			services.GET("", middleware.HasPermission("services", "read"), serviceHandler.ListServices)
			services.GET("/:id", middleware.HasPermission("services", "read"), serviceHandler.GetService)

			// Modify services - requires services:write permission
			services.POST("", middleware.HasPermission("services", "write"), serviceHandler.CreateService)
			services.PUT("/:id", middleware.HasPermission("services", "write"), serviceHandler.UpdateService)
			services.DELETE("/:id", middleware.HasPermission("services", "write"), serviceHandler.DeleteService)
		}
	}
}
