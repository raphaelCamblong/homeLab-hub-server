package seed

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

// SeedPermissions creates all CRUD permissions for the app's resources
func SeedPermissions(db database.Database) {
	if db.Get().Migrator().HasTable(&entities.Permission{}) {
		logrus.Info("✅ Permissions table already exists")
		return
	}

	verbs := []string{"create", "read", "update", "delete", "execute"}
	resources := []string{
		"users",
		"roles",
		"permissions",
		"api-keys",
		"audit",
		"services",
		"pipelines",
		"infrastructure",
	}

	for _, resource := range resources {
		for _, verb := range verbs {
			name := fmt.Sprintf("%s:%s", resource, verb)
			permission := entities.Permission{
				Name:        name,
				Description: fmt.Sprintf("Can %s %s", verb, resource),
				Resource:    resource,
				Action:      verb,
			}

			err := db.Get().FirstOrCreate(&permission, entities.Permission{Name: name}).Error
			if err != nil {
				logrus.Errorf("❌ Failed to create permission %s: %v", name, err)
			} else {
				logrus.Infof("✅ Permission seeded: %s", name)
			}
		}
	}
}

// SeedRoles seeds core roles and associates full permissions to admin
func SeedRoles(db database.Database) {
	if db.Get().Migrator().HasTable(&entities.Role{}) {
		logrus.Info("✅ Roles table already exists")
		return
	}

	roles := []entities.Role{
		{Name: "admin", Description: "Administrator role with full access"},
		{Name: "user", Description: "Regular user role with limited access"},
		{Name: "guest", Description: "Guest role with limited access"},
	}

	for _, role := range roles {
		if err := db.Get().FirstOrCreate(&role, entities.Role{Name: role.Name}).Error; err != nil {
			logrus.Errorf("❌ Failed to create role %s: %v", role.Name, err)
		} else {
			logrus.Infof("✅ Role seeded: %s", role.Name)
		}
	}

	var adminRole entities.Role
	var allPermissions []entities.Permission

	db.Get().Where("name = ?", "admin").First(&adminRole)
	db.Get().Find(&allPermissions)

	if err := db.Get().Model(&adminRole).Association("Permissions").Replace(allPermissions); err != nil {
		logrus.Errorf("❌ Failed to attach permissions to admin role: %v", err)
	} else {
		logrus.Infof("✅ Admin role now has %d permissions", len(allPermissions))
	}
}
