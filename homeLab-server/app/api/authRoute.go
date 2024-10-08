package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func AuthRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1")

	handler := handlers.NewAuthenticationHandler(usecase.NewAuthenticationUseCase(repo.Auth))

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	return nil
}
