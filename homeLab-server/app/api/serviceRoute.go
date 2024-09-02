package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func ServiceRoute(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/service")
	router.Use(middleware.JWTAuthMiddleware())

	handler := handlers.NewServiceHandler(usecase.NewServiceUseCase(repo.Service))

	router.GET("/services", handler.GetAllService)
	router.GET("/service/:id", handler.GetServiceById)

	return nil
}
