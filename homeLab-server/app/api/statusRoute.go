package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func StatusRoutes(infra *infrastructure.Infrastructure) error {
	router := infra.Router.Get()
	db := infra.Db

	statusHandler := handlers.NewStatusHandler(usecase.NewStatusUseCase(repositories.NewStatusRepository(db)))

	router.Group("").GET("/status", statusHandler.GetStatus)
	return nil
}
