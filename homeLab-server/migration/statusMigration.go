package migration

import (
	"log"

	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

func StatusMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(&entities.StatusEntity{})
	if err != nil {
		log.Fatal(err)
	}
	db.GetDb().CreateInBatches(
		[]*entities.StatusEntity{
			{
				Version: "0.001",
				Health:  "critical",
			},
			{Version: "0.002", Health: "critical"},
			{Version: "0.003", Health: "critical"},
		}, 1,
	)
}
