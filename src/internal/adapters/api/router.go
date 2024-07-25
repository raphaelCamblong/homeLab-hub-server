package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	r := gin.Default()
	s := &Server{
		engine:  r,
		handler: handler,
	}

	r.POST("/services", s.handler.CreateService)

	return s
}

func (s *Server) Start() error {
	return s.engine.Run(":8080") // Replace ":8080" with your desired port
}
