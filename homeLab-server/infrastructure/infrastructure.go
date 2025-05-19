package infrastructure

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
	"homelab.com/homelab-server/homeLab-server/infrastructure/streaming"
)

type Infrastructure struct {
	Router              router.Router
	Db                  database.Database
	ExternalHttpService client.ExternalHttpService
	StreamHub           *streaming.StreamHub
	Cron                *cron.Cron
}
