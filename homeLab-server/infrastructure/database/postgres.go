package database

import (
	"fmt"
	"log"

	"gorm.io/gorm/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func NewPostgresDatabase(c config.DatabaseConfig) *database {
	logrus.Debugf("Database connection parameters: Host=%s, User=%s, Port=%d, DBName=%s",
		c.Host, c.User, c.Port, c.DBName)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", c.Host, c.User, c.Password, c.DBName, c.Port)
	logrus.Debugf("Constructed DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	logrus.Info("Connected to:", dsn)
	return &database{Db: db}
}
