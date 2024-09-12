package router

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"
	"sync"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/app/config"
)

type GinRouter struct {
	Router *gin.Engine
}

var (
	once           sync.Once
	routerInstance *GinRouter
)

func NewRouter() (Router, error) {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"192.168.1.0/24", "10.0.0.0/8"})
	r.Use(middleware.CORSMiddleware())
	once.Do(
		func() {
			routerInstance = &GinRouter{
				Router: r,
			}
		},
	)
	return routerInstance, nil
}

func (s *GinRouter) Start() {
	c := config.GetConfig()
	addr := fmt.Sprintf(":%d", c.Server.Port)

	err := s.Router.Run(addr)
	if err != nil {
		logrus.Errorf("failed to iniate gin router %d", err)
	}
}

func (s *GinRouter) Get() *gin.Engine {
	return routerInstance.Router
}
