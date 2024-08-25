package app

import (
	"homelab.com/homelab-server/homeLab-server/app/api"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/externalHttpService"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
)

type App struct {
	Infra *infrastructure.Infrastructure
}

func (app *App) Start() {
	_ = api.StatusRoutes(app.Infra)
	_ = api.HealthRoute(app.Infra)
	_ = api.ThermalRoutes(app.Infra)
	app.Infra.Router.Start()
}

func NewApp() *App {
	return &App{
		Infra: &infrastructure.Infrastructure{
			Router:              initRouter(),
			Cache:               initRedis(),
			Db:                  initSQLite(),
			ExternalHttpService: initExternalHttpService(),
		},
	}
}

func initRedis() cache.Database {
	redisDb, err := cache.NewRedisDatabase()
	if err != nil {
		panic(err)
	}
	defer func() { _ = redisDb.GetClient().Close() }()
	return redisDb
}

func initSQLite() database.Database {
	db, err := database.NewSqliteDatabase()
	if err != nil {
		panic(err)
	}
	return db
}

func initRouter() router.Router {
	rtr, err := router.NewRouter()
	if err != nil {
		panic(err)
	}

	return rtr
}

func initExternalHttpService() externalHttpService.ExternalHttpService {
	cfg := config.GetConfig()
	redfish := externalHttpService.NewRedfishInfra(cfg.ExternalServicesCredential.IloIp)
	return externalHttpService.NewExternalHttpService(redfish, nil)
}
