package main

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/cache"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router"
	"log"
)

func initRedis() cache.Database {
	redisDb, err := cache.NewRedisDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer redisDb.GetClient().Close()
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

type App struct {
	router router.Router
	cache  cache.Database
	db     database.Database
}

func main() {
	app := App{
		router: initRouter(),
		cache:  initRedis(),
		db:     initSQLite(),
	}

	app.router.Start()
}
