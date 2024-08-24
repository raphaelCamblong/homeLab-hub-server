package api

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"net/http"
)

func HealthRoute(infra *infrastructure.Infrastructure) error {
	r := infra.Router.Get()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server Ok"})
	})
	return nil
}
