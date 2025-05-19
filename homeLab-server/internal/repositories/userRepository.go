package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/pkg/tools"
)

type UserRepository interface {
	// Authentication
	Login(credentials entities.UserCredentials, ipAddress, userAgent string) (*entities.AuthResponse, error)
	Logout(userID uint) error
	Register(credentials entities.UserCredentials) (*entities.AuthResponse, error)
	GetMe(userID uint) (*entities.UserResponse, error)
	ResetPassword(userID uint, newPassword string) error

	// User Management
	GetUser(id uint) (*entities.UserEntity, error)
	GetUserByUsername(username string) (*entities.UserEntity, error)
	UpdateUser(id uint, user entities.UserEntity) error
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]entities.UserResponse, int64, error)
	SetUserActive(id uint, active bool) error

	// Role Management
	AssignRole(userID uint, roleID uint) error
	RemoveRole(userID uint, roleID uint) error
	CreateRole(role entities.Role) error
	UpdateRole(id uint, role entities.Role) error
	DeleteRole(id uint) error
	ListRoles() ([]entities.Role, error)

	// Permission Management
	ListPermissions() ([]entities.Permission, error)
	GetUserPermissions(userID uint) ([]entities.Permission, error)
	AssignPermission(roleID uint, permissionID uint) error
	RemovePermission(roleID uint, permissionID uint) error

	// API Key Management
	CreateAPIKey(userID uint, name string, permissions []string, expiresAt *time.Time) (*entities.APIKey, error)
	RevokeAPIKey(id uint) error
	ListAPIKeys(userID uint) ([]entities.APIKey, error)

	// Audit
	ViewAuditLogs(userID uint, page, pageSize int) ([]entities.AuditLog, int64, error)
	LogAudit(log entities.AuditLog) error
}

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Login(cred entities.UserCredentials, ipAddress, userAgent string) (*entities.AuthResponse, error) {
	var user entities.UserEntity
	if err := r.db.Get().Preload("Roles").Where("username = ?", cred.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if user.IsLocked() {
		return nil, errors.New("account is temporarily locked")
	}

	if !tools.CheckPasswordHash(cred.Password, user.Password) {
		user.LoginAttempts++
		if user.LoginAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lockUntil
		}
		r.db.Get().Save(&user)
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"last_login_at":   &now,
		"last_ip_address": &ipAddress,
		"user_agent":      &userAgent,
		"login_attempts":  0,
		"locked_until":    nil,
	}
	r.db.Get().Model(&user).Updates(updates)

	// Generate roles for token
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Generate token
	cfg := config.Get().Core.API.Security.JWT
	token, err := tools.GenerateJWT(user.Username, user.ID, cfg.Expiration, cfg.Secret)
	if err != nil {
		return nil, err
	}

	// Log successful login
	r.LogAudit(entities.AuditLog{
		UserID:    user.ID,
		Action:    "login",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	jwtToken := entities.JwtToken(token)
	response := entities.AuthResponse{
		Token: jwtToken,
		User:  user.ToResponse(),
	}
	return &response, nil
}

func (r *userRepository) Register(cred entities.UserCredentials) (*entities.AuthResponse, error) {
	var existingUser entities.UserEntity
	if err := r.db.Get().Where("username = ?", cred.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := tools.HashPassword(cred.Password)
	if err != nil {
		return nil, err
	}

	newUser := entities.UserEntity{
		Username: cred.Username,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := r.db.Get().Create(&newUser).Error; err != nil {
		return nil, err
	}

	// Assign default role if exists
	var defaultRole entities.Role
	if err := r.db.Get().Where("name = ?", "user").First(&defaultRole).Error; err == nil {
		r.AssignRole(newUser.ID, defaultRole.ID)
	}

	return r.Login(cred, "", "")
}

func (r *userRepository) AssignRole(userID uint, roleID uint) error {
	return r.db.Get().Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
}

func (r *userRepository) RemoveRole(userID uint, roleID uint) error {
	return r.db.Get().Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Error
}

func (r *userRepository) LogAudit(log entities.AuditLog) error {
	return r.db.Get().Create(&log).Error
}

func (r *userRepository) Logout(userID uint) error {
	// Log the logout action
	return r.LogAudit(entities.AuditLog{
		UserID: userID,
		Action: "logout",
	})
}

func (r *userRepository) ResetPassword(userID uint, newPassword string) error {
	hashedPassword, err := tools.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return r.db.Get().Model(&entities.UserEntity{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *userRepository) GetUser(id uint) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := r.db.Get().Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := r.db.Get().Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(id uint, user entities.UserEntity) error {
	return r.db.Get().Model(&entities.UserEntity{}).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Get().Delete(&entities.UserEntity{}, id).Error
}

func (r *userRepository) ListUsers(page, pageSize int) ([]entities.UserResponse, int64, error) {
	var users []entities.UserEntity
	var total int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.Get().Model(&entities.UserEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated users
	if err := r.db.Get().Preload("Roles").Preload("Roles.Permissions").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	userResponses := make([]entities.UserResponse, len(users))
	for i, user := range users {
		roles := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roles[j] = role.Name
		}

		userResponses[i] = entities.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			IsActive:  user.IsActive,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
		}
	}

	return userResponses, total, nil
}

func (r *userRepository) SetUserActive(id uint, active bool) error {
	return r.db.Get().Model(&entities.UserEntity{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *userRepository) CreateRole(role entities.Role) error {
	return r.db.Get().Create(&role).Error
}

func (r *userRepository) UpdateRole(id uint, role entities.Role) error {
	return r.db.Get().Model(&entities.Role{}).Where("id = ?", id).Updates(role).Error
}

func (r *userRepository) DeleteRole(id uint) error {
	return r.db.Get().Delete(&entities.Role{}, id).Error
}

func (r *userRepository) ListRoles() ([]entities.Role, error) {
	var roles []entities.Role
	err := r.db.Get().Find(&roles).Error
	return roles, err
}

func (r *userRepository) ListPermissions() ([]entities.Permission, error) {
	var permissions []entities.Permission
	err := r.db.Get().Find(&permissions).Error
	return permissions, err
}

func (r *userRepository) AssignPermission(roleID uint, permissionID uint) error {
	return r.db.Get().Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error
}

func (r *userRepository) RemovePermission(roleID uint, permissionID uint) error {
	return r.db.Get().Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID).Error
}

func (r *userRepository) CreateAPIKey(userID uint, name string, permissions []string, expiresAt *time.Time) (*entities.APIKey, error) {
	apiKeyString := tools.GenerateRandomString(32)

	apiKey := &entities.APIKey{
		UserID:      userID,
		Name:        name,
		Key:         apiKeyString,
		Permissions: permissions,
		ExpiresAt:   expiresAt,
	}

	if err := r.db.Get().Create(apiKey).Error; err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (r *userRepository) RevokeAPIKey(id uint) error {
	return r.db.Get().Delete(&entities.APIKey{}, id).Error
}

func (r *userRepository) ListAPIKeys(userID uint) ([]entities.APIKey, error) {
	var apiKeys []entities.APIKey
	err := r.db.Get().Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

func (r *userRepository) ViewAuditLogs(userID uint, page, pageSize int) ([]entities.AuditLog, int64, error) {
	var logs []entities.AuditLog
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Get().Model(&entities.AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Get().Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *userRepository) GetMe(userID uint) (*entities.UserResponse, error) {
	var user entities.UserEntity
	if err := r.db.Get().Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (r *userRepository) GetUserPermissions(userID uint) ([]entities.Permission, error) {
	var permissions []entities.Permission

	err := r.db.Get().
		Select("DISTINCT permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error

	return permissions, err
}
