package database

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func NewSqliteDatabase(c config.DatabaseConfig) *database {
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.Host)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("failed to connect database: %d", err)
	}
	dbInstance = &database{Db: db}
	logrus.Infof("Successfully connected to database: %s", dsn)
	return dbInstance
}
