package main

import (
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/cmd/migrate/seed"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/init/config"
)

func main() {
	c := config.Get()
	dotEnvConfig := config.LoadDotEnv(".env")
	if err := c.FlattenConfig(*dotEnvConfig); err != nil {
		logrus.Error("❌ Failed to flatten config: %w", err)
		panic(err)
	}
	db, err := database.NewDatabase(c.Client.Database)
	if err != nil {
		logrus.Error("Failed to connect to database: ", err)
		return
	}

	AutoMigrate(db)
}

func AutoMigrate(db database.Database) {
	seed.SeedPipelines(db)

	logrus.Info("✅ Migration and seeding complete.")
}
