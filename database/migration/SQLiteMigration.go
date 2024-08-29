package main

import (
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/migration"
)

func main() {
	db, _ := database.NewSqliteDatabase()
	migration.StatusMigration(db)
	logrus.Info("Successfully Run the migration!")
}
