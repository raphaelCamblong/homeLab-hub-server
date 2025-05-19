package main

import (
	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/cmd/migrate/seed"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
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
	logrus.Info("Migrating database...")
	logrus.Info("Running AutoMigrate...")
	db.Get().AutoMigrate(
		&entities.UserEntity{},
		&entities.AuditLog{},
		&entities.Role{},
		&entities.Permission{},
		&entities.APIKey{},
		&entities.ServiceEntity{},
		&entities.InfrastructureEntity{},
		&entities.PipelineTemplate{},
		&entities.StepTemplate{},
		&entities.Job{},
		&entities.Step{},
		&entities.ServiceEntity{},
	)

	logrus.Info("Done.")

	// Seeding logic
	// seed.SeedPermissions(db)
	// seed.SeedRoles(db)
	// seed.SeedAdminUser(db)
	seed.SeedPipelines(db)
	// seed.SeedServices(db)

	logrus.Info("✅ Migration and seeding complete.")
}
