package app

import (
	"homelab.com/homelab-server/homeLab-server/app/api"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
	"log"
)

type App struct {
	Infra *infrastructure.Infrastructure
}

func (app *App) Start() {
	_ = api.StatusRoutes(app.Infra)
	_ = api.HealthRoute(app.Infra)
	app.Infra.Router.Start()
}

func NewApp() *App {
	return &App{
		Infra: &infrastructure.Infrastructure{
			Router:              initRouter(),
			Cache:               nil, //initRedis()
			Db:                  initSQLite(),
			ExternalHttpService: nil,
		},
	}
}

func initRedis() cache.Database {
	redisDb, err := cache.NewRedisDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = redisDb.GetClient().Close() }()
	return redisDb
}

func initSQLite() database.Database {
	db, err := database.NewSqliteDatabase()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initRouter() router.Router {
	rtr, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	return rtr
}
