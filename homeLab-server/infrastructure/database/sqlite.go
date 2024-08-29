package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/app/config"
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
	once.Do(
		func() {
			dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.Db.Host)

			db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
			if err != nil {
				logrus.Errorf("failed to connect database: %d", err)
			}
			dbInstance = &sqliteDatabase{Db: db}
			logrus.Infof("Successfully connected to database: %s", dsn)
		},
	)
	return dbInstance, nil
}

func (s *sqliteDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
