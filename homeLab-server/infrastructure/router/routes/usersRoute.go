package routes

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func UserRoutes(infra *infrastructure.Infrastructure, userRepo repositories.UserRepository) {
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Public
	public := infra.Router.Get().Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
		}
	}

	// Protected
	protected := infra.Router.Get().Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		auth := protected.Group("/auth")
		{
			auth.POST("/logout", userHandler.Logout)
		}

		users := protected.Group("/users")
		{
			// Get current user
			users.GET("/me", userHandler.GetMe)

			// List and get users - requires users:read permission
			users.GET("", middleware.HasPermission("users", "read"), userHandler.ListUsers)
			users.GET("/:id", middleware.HasPermission("users", "read"), userHandler.GetUser)

			// Role management - requires roles:read/write permission
			users.GET("/roles", middleware.HasPermission("roles", "read"), userHandler.ListRoles)
			users.GET("/:id/roles", middleware.HasPermission("roles", "read"), userHandler.GetUserRoles)
			users.POST("/:id/roles/:roleId", middleware.HasPermission("roles", "write"), userHandler.AssignRole)
			users.DELETE("/:id/roles/:roleId", middleware.HasPermission("roles", "write"), userHandler.RemoveRole)

			// Permission management - requires permissions:read/write permission
			users.GET("/permissions", middleware.HasPermission("permissions", "read"), userHandler.ListPermissions)
			users.GET("/:id/permissions", middleware.HasPermission("roles", "read"), userHandler.GetUserPermissions)
			users.POST("/:id/permissions/:permissionId", middleware.HasPermission("roles", "write"), userHandler.AssignPermission)
			users.DELETE("/:id/permissions/:permissionId", middleware.HasPermission("roles", "write"), userHandler.RemovePermission)
		}

		apiKeys := users.Group("/api-keys")
		{
			// API Key management - only requires authentication
			apiKeys.POST("", userHandler.CreateAPIKey)
			apiKeys.GET("", userHandler.ListAPIKeys)
			apiKeys.DELETE("/:keyId", userHandler.RevokeAPIKey)
		}

		protected.GET("/audit-logs", middleware.HasPermission("audit", "read"), userHandler.GetAuditLogs)
	}
}
