package migration

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"os"
)

func loadJson() (*[]entities.ServiceEntity, error) {
	file, err := os.Open("./database/migration/service.json")
	if err != nil {
		logrus.Error("Failed to open JSON file: ", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var services []entities.ServiceEntity
	err = decoder.Decode(&services)
	if err != nil {
		logrus.Error("Failed to parse JSON file: ", err)
		return nil, err
	}

	return &services, nil
}

func ServiceMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(&entities.ServiceEntity{})
	if err != nil {
		logrus.Error(err)
		return
	}

	services, err := loadJson()
	if err != nil {
		logrus.Error("Failed to load JSON file: ", err)
		return
	}
	db.GetDb().CreateInBatches(services, 10)
}
