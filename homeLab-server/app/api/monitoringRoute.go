package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func MonitoringRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/monitoring")

	handler := handlers.NewMonitoringHandler(usecase.NewMonitoringUseCase(repo.Monitoring))

	router.GET("/prometheus", handler.QueryPrometheus)
	router.GET("/grafana", handler.GetGrafanaDashboardData)
	router.GET("/alertmanager", handler.FetchActiveAlerts)

	return nil
} 