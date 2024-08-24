package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/app/config"
	"log"
	"sync"
)

type sqliteDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *sqliteDatabase
)

func NewSqliteDatabase() (Database, error) {
	c := config.GetConfig()
	once.Do(func() {
		dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.Db.Host)

		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}
		dbInstance = &sqliteDatabase{Db: db}
		log.Default().Print("Successfully connected to database: %s", dsn)
	})
	return dbInstance, nil
}

func (s *sqliteDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
