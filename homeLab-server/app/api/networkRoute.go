package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func NetworkRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/network")

	handler := handlers.NewNetworkHandler(usecase.NewNetworkUseCase(repo.Network))

	router.GET("/devices", handler.ListDevices)
	router.GET("/interfaces", handler.ListNetworkInterfaces)
	router.GET("/interfaces/{id}", handler.GetInterfaceDetails)
	router.GET("/firewall", handler.ViewFirewallRules)
	router.POST("/firewall/rules", handler.ModifyFirewallRules)
	router.GET("/dns", handler.ViewDNSSettings)
	router.POST("/speedtest", handler.RunNetworkSpeedTest)

	return nil
} 