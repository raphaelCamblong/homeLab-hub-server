package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/pkg/tools"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := tools.ParseJWT(tokenStr, config.Get().Core.API.Security.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("user_id", strconv.FormatUint(uint64(claims.UserID), 10))
		if claims.Roles != nil {
			c.Set("roles", claims.Roles)
		}
		SetUserContext(c)
		c.Next()
	}
}

func SetUserContext(c *gin.Context) {
	dbAny, exists := c.Get("mainDb")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not found in context"})
		c.Abort()
		return
	}
	db := dbAny.(database.Database)
	userID := c.GetString("user_id")

	var user entities.UserEntity
	if err := db.Get().
		Preload("Roles").
		Preload("Roles.Permissions").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
		c.Abort()
		return
	}
	c.Set("user", &user)
}
