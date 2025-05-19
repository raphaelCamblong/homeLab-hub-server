package tools

import (
	"fmt"

	"homelab.com/homelab-server/homeLab-server/docs"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func UpdateDocs(config *config.Config) {
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", config.Core.API.Host, config.Core.API.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Title = "HomeLab API"
	docs.SwaggerInfo.Description = "HomeLab API specification"
}
