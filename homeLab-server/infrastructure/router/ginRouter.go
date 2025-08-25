package router

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/middleware"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

type GinRouter struct {
	Router *gin.Engine
}

var (
	once           sync.Once
	routerInstance *GinRouter
)

func NewRouter() (Router, error) {
	c := config.Get()
	r := gin.Default()
	_ = r.SetTrustedProxies(c.Core.API.Security.TrustedProxies)
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
	c := config.Get()
	addr := fmt.Sprintf("%s:%d", c.Core.API.Host, c.Core.API.Port)

	err := s.Router.Run(addr)
	if err != nil {
		logrus.Errorf("failed to iniate gin router %d", err)
	}
	logrus.Infof("gin router started on %s", addr)
}

func (s *GinRouter) Get() *gin.Engine {
	return routerInstance.Router
}
