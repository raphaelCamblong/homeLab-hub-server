package processor

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func initDB() *database.Database {
	config := config.Get()
	db, err := database.NewDatabase(config.Client.Database)
	if err != nil {
		panic(err)
	}
	return &db
}

func initRouter() *router.Router {
	rtr, err := router.NewRouter()

	if err != nil {
		panic(err)
	}
	return &rtr
}

func initExternalHttpService() client.ExternalHttpService {
	return client.NewExternalHttpService()
}
