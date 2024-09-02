package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func IloRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/ilo")
	router.Use(middleware.JWTAuthMiddleware())

	thermalHandler := handlers.NewIloHandler(usecase.NewRewIloUseCase(repo.Ilo))

	router.GET("/thermal", thermalHandler.GetThermal)
	router.GET("/power", thermalHandler.GetPower)
	return nil
}
