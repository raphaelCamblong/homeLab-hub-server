package migration

import (
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
)

func UserMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(&entities.UserEntity{})
	if err != nil {
		logrus.Error(err)
	}
	db.GetDb().CreateInBatches(
		[]*entities.UserEntity{
			{
				Username: "test1",
				Password: "invalid",
			},
		}, 1,
	)
}
