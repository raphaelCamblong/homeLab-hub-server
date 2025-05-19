package database

import (
	"sync"

	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

type Database interface {
	Get() *gorm.DB
}

type database struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *database
)

func (d *database) Get() *gorm.DB {
	return d.Db
}

func NewDatabase(config config.DatabaseConfig) (Database, error) {
	once.Do(
		func() {
			dbInstance = NewPostgresDatabase(config)
		},
	)
	return dbInstance, nil
}
