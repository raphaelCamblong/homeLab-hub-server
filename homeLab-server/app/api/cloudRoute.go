package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func CloudRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/cloud")
	router.Use(middleware.JWTAuthMiddleware())

	handler := handlers.NewCloudHandler(usecase.NewCloudUseCase(repo.Cloud))

	router.GET("/vms", handler.GetVmsData)
	router.GET("/host", handler.GetHostData)
	return nil
}
