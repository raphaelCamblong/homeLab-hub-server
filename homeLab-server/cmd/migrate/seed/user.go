package seed

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

func SeedAdminUser(db database.Database) {
	var adminUser entities.UserEntity
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	// Check if user exists
	if err := db.Get().FirstOrCreate(&adminUser, entities.UserEntity{Username: "admin"}).Error; err != nil {
		logrus.Errorf("❌ Failed to create admin user: %v", err)
		return
	}

	// Update password if it's empty
	if adminUser.Password == "" {
		adminUser.Password = string(hashedPassword)
		db.Get().Save(&adminUser)
	}

	// Assign admin role
	var adminRole entities.Role
	if err := db.Get().Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		logrus.Errorf("❌ Admin role not found: %v", err)
		return
	}

	if err := db.Get().Model(&adminUser).Association("Roles").Replace(&adminRole); err != nil {
		logrus.Errorf("❌ Failed to assign admin role to user: %v", err)
	} else {
		logrus.Infof("✅ Admin user seeded and assigned to admin role")
	}
}
