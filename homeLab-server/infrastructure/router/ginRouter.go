package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/config"
	"net/http"
)

type GinRouter struct {
	router        *gin.Engine
	publicRouter  *gin.RouterGroup
	privateRouter *gin.RouterGroup
}

func NewRouter() (Router, error) {
	r := gin.Default()
	publicGroup := r.Group("")
	privateGroup := r.Group("")
	s := &GinRouter{
		router:        r,
		publicRouter:  publicGroup,
		privateRouter: privateGroup,
	}
	return s, nil
}

func (s *GinRouter) Start() {
	c := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": ""})
	})

	err := s.router.Run(addr)
	if err != nil {
		panic(err)
	}
}
