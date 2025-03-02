package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func UserRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/users")

	handler := handlers.NewUserHandler(usecase.NewUserUseCase(repo.Users))

	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)
	router.GET("/permissions", handler.ViewUserRoles)
	router.GET("/audit-logs", handler.ViewAuditLogs)

	return nil
} 