package main

import (
	"log"

	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/migration"
)

func main() {
	db, _ := database.NewSqliteDatabase()
	migration.StatusMigration(db)
	log.Print("Successfully Run the migration!")
}
