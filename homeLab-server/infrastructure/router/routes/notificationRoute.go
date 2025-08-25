package routes

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func NotificationRoutes(infra *infrastructure.Infrastructure, notificationRepo repositories.NotificationRepository) {
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)

	protected := infra.Router.Get().Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		notifications := protected.Group("/notifications")
		{
			// List notifications
			notifications.GET("", notificationHandler.ListNotifications)

			// Stream notifications
			notifications.GET("/stream", notificationHandler.StreamNotifications)
		}
	}
}
