package router

import "github.com/gin-gonic/gin"

type Router interface {
	Start()
	Get() *gin.Engine
}
