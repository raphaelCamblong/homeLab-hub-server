package app

import (
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/app/api"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/externalHttpService"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
)

type (
	App struct {
		Infra        *infrastructure.Infrastructure
		Repositories *api.Repositories
	}
)

func NewApp() *App {
	return &App{
		Infra: &infrastructure.Infrastructure{
			Router:              *initRouter(),
			Cache:               nil, //*initRedis(),
			Db:                  *initSQLite(),
			ExternalHttpService: initExternalHttpService(),
		},
	}
}

func (app *App) Start() {
	logrus.Info("Loading repositories")
	app.Repositories = api.NewRepositories(app.Infra)

	logrus.Info("Binding api routes")
	_ = api.HealthRoute(app.Infra, app.Repositories)
	_ = api.AuthRoutes(app.Infra, app.Repositories)
	_ = api.IloRoutes(app.Infra, app.Repositories)
	_ = api.CloudRoutes(app.Infra, app.Repositories)
	_ = api.ServiceRoute(app.Infra, app.Repositories)

	logrus.Info("Starting router")
	app.Infra.Router.Start()
}

func initRedis() *cache.Database {
	redisDb, err := cache.NewRedisDatabase()
	if err != nil {
		logrus.Error("Failed to connect to Redis: ", err)
		//panic(err)
		return nil
	}
	return &redisDb
}

func initSQLite() *database.Database {
	db, err := database.NewSqliteDatabase()
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

func initExternalHttpService() externalHttpService.ExternalHttpService {
	cfg := config.GetConfig()
	redfish := externalHttpService.NewRedfishInfra(cfg.ExternalServicesCredential.Ilo)
	xo := externalHttpService.NewXenOrchestraInfra(cfg.ExternalServicesCredential.XO)
	return externalHttpService.NewExternalHttpService(redfish, xo)
}
