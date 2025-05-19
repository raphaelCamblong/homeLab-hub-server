package routes

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
)

func DocsRoutes(infra *infrastructure.Infrastructure) {
	infra.Router.Get().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
