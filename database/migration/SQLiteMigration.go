package main

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/migration"
	"log"
)

func main() {
	db, _ := database.NewSqliteDatabase()
	migration.StatusMigration(db)
	log.Print("Successfully Run the migration!")
}
