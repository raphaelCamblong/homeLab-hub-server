package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func ClusterRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/Clusters")

	handler := handlers.NewClusterHandler(usecase.NewClusterUseCase(repo.Cluster))

	router.GET("/nodes", handler.ListClusterNodes)
	router.GET("/pods", handler.ListRunningPods)
	router.GET("/pods/{namespace}", handler.ListPodsInNamespace)
	router.GET("/services", handler.ListClusterServices)
	router.GET("/deployments", handler.ListClusterDeployments)
	router.GET("/logs/{pod}", handler.GetPodLogs)
	router.GET("/events", handler.GetClusterEvents)
	router.POST("/apply", handler.ApplyYAMLManifest)
	router.DELETE("/delete", handler.DeleteResource)

	return nil
} 