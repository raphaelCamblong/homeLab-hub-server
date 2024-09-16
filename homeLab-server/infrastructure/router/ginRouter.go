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
	c := config.GetConfig()
	r := gin.Default()
	_ = r.SetTrustedProxies(c.App.Security.TrustedProxies)
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
	addr := fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)

	err := s.Router.Run(addr)
	if err != nil {
		logrus.Errorf("failed to iniate gin router %d", err)
	}
}

func (s *GinRouter) Get() *gin.Engine {
	return routerInstance.Router
}
