package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func SystemRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/system")

	handler := handlers.NewSystemHandler(usecase.NewSystemUseCase(repo.System))

	router.GET("/info", handler.GetSystemInfo)
	router.GET("/processes", handler.ListProcesses)
	router.GET("/services", handler.ListActiveServices)
	router.GET("/metrics", handler.FetchSystemMetrics)
	router.GET("/health", handler.CheckHealth)
	router.POST("/restart", handler.RestartSystem) // restricted
	router.POST("/shutdown", handler.ShutdownSystem) // restricted

	return nil
} 