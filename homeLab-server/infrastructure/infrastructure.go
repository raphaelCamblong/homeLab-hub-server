package infrastructure

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
)

type Infrastructure struct {
	Router              router.Router
	Cache               cache.Database
	Db                  database.Database
	ExternalHttpService any
}
