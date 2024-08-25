package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func ThermalRoutes(infra *infrastructure.Infrastructure) error {
	router := infra.Router.Get()

	redfishRepo := repositories.NewRedfishRepository(infra.Cache, infra.ExternalHttpService.GetRedfish())
	thermalRepo := repositories.NewThermalRepository(redfishRepo)

	thermalHandler := handlers.NewThermalHandler(usecase.NewThermalUseCase(thermalRepo))

	router.Group("").GET("/thermal", thermalHandler.GetThermal)
	return nil
}
