package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func HealthRoute(infra *infrastructure.Infrastructure, repo *Repositories) error {
	r := infra.Router.Get().Group("/api/v1")

	statusHandler := handlers.NewHealthHandler(usecase.NewStatusUseCase(repo.Status))

	infra.Router.Get().GET("/status", statusHandler.GetOK)
	r.GET("/status", statusHandler.GetStatus)
	r.GET("/health", statusHandler.GetOK)

	return nil
}
