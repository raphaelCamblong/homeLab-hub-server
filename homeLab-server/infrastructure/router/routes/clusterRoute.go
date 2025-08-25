package routes

import (
	"homelab.com/homelab-server/homeLab-server/internal/handlers"

	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func ClusterRoutes(infra *infrastructure.Infrastructure, clusterRepo repositories.ClusterRepository) {
	clusterUseCase := usecase.NewClusterUseCase(clusterRepo)
	clusterHandler := handlers.NewClusterHandler(clusterUseCase)

	v1 := infra.Router.Get().Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		cluster := v1.Group("/cluster")
		{
			cluster.GET("/", clusterHandler.GetClusterStats)
			cluster.GET("/nodes", clusterHandler.GetClusterNodes)
			cluster.GET("/nodes/:nodeID/stats", clusterHandler.GetNodeStats)
			cluster.GET("/pods", clusterHandler.GetRunningPods)
			cluster.GET("/pods/:namespace", clusterHandler.GetPodsInNamespace)
			cluster.GET("/workflows", clusterHandler.GetWorkloads)
			cluster.GET("/services", clusterHandler.GetClusterServices)
			cluster.GET("/deployments", clusterHandler.GetClusterDeployments)
		}
	}
}
