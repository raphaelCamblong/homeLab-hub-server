package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func K3sRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/k3s")

	handler := handlers.NewK3sHandler(usecase.NewK3sUseCase(repo.K3s))

	router.GET("/nodes", handler.ListK3sNodes)
	router.GET("/pods", handler.ListRunningPods)
	router.GET("/pods/{namespace}", handler.ListPodsInNamespace)
	router.GET("/services", handler.ListK3sServices)
	router.GET("/deployments", handler.ListK3sDeployments)
	router.GET("/logs/{pod}", handler.GetPodLogs)
	router.GET("/events", handler.GetClusterEvents)
	router.POST("/apply", handler.ApplyYAMLManifest)
	router.DELETE("/delete", handler.DeleteResource)

	return nil
} 