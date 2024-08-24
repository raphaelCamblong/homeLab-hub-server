package router

import (
	"fmt"
	"log"
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
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)

	err := s.Router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *GinRouter) Get() *gin.Engine {
	return routerInstance.Router
}
