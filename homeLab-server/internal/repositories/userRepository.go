package repositories

import (
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"errors"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/cmd/tools"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
)

type UserRepository interface {
	Login(credentials entities.UserCredentials) (*entities.JwtToken, error)
	Logout(userID string) error
	ViewUserRoles() ([]entities.Role, error)
	ViewAuditLogs() ([]entities.AuditLog, error)
	Register(credentials entities.UserCredentials) (*entities.JwtToken, error)
}

type userRepository struct {
	db database.Database
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Logout(userID string) error {
	// Implement the logic to log out a user
	// For example, invalidate the user's session or token
	return nil
}

func (r *userRepository) ViewUserRoles() ([]entities.Role, error) {
	// Implement the logic to view user roles
	return []entities.Role{}, nil // Placeholder return
}

func (r *userRepository) ViewAuditLogs() ([]entities.AuditLog, error) {
	var auditLogs []entities.AuditLog
	if err := r.db.GetDb().Find(&auditLogs).Error; err != nil {
		return nil, errors.New("failed to retrieve audit logs: " + err.Error())
	}
	return auditLogs, nil
}

func (r *userRepository) Register(cred entities.UserCredentials) (*entities.JwtToken, error) {
	var existingUser entities.UserEntity
	if err := r.db.GetDb().Where("username = ?", cred.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := tools.HashPassword(cred.Password)
	if err != nil {
		return nil, err
	}

	newUser := entities.UserEntity{
		Username: cred.Username,
		Password: hashedPassword,
	}

	if err := r.db.GetDb().Create(&newUser).Error; err != nil {
		return nil, err
	}

	return r.Login(cred)
}

func (r *userRepository) Login(cred entities.UserCredentials) (*entities.JwtToken, error) {
	var user entities.UserEntity
	if err := r.db.GetDb().Where("username = ?", cred.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !tools.CheckPasswordHash(cred.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := tools.GenerateJWT(cred.Username)
	if err != nil {
		return nil, err
	}

	jwtToken := entities.JwtToken(token)
	return &jwtToken, nil
}
