package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func NewPostgresDatabase(c config.DatabaseConfig) *database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", c.Host, c.User, c.Password, c.DBName, c.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	logrus.Info("Connected to:", dsn)
	return &database{Db: db}
}
