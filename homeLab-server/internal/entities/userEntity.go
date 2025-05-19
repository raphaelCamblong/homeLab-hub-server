package entities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model

	// Core user information
	Username    string     `gorm:"uniqueIndex;size:255" json:"username"`
	Password    string     `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// Role-based access control
	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`

	// Security
	LoginAttempts int        `json:"-" gorm:"default:0"`
	LockedUntil   *time.Time `json:"-"`

	// API Access
	APIKeys []APIKey `gorm:"foreignKey:UserID" json:"-"`

	// Audit
	LastIPAddress *string `json:"last_ip_address"`
	UserAgent     *string `json:"user_agent"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"uniqueIndex" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex" json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"` // create, read, update, delete
}

type APIKey struct {
	gorm.Model
	UserID      uint       `json:"user_id"`
	Name        string     `json:"name"`
	Key         string     `json:"-"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Permissions []string   `gorm:"type:json" json:"permissions"`
}

type AuditLog struct {
	gorm.Model
	UserID     uint                   `json:"user_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Details    map[string]interface{} `gorm:"type:json" json:"details"`
}

// Authentication related structs
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtToken string

type Claims struct {
	Username string   `json:"username"`
	UserID   uint     `json:"user_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	Token JwtToken     `json:"token"`
	User  UserResponse `json:"user"`
}

// Response structs
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Roles     []string  `json:"roles"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at" swaggertype:"primitive,string" `
}

func (u *UserEntity) ToResponse() UserResponse {
	roles := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = role.Name
	}

	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Roles:     roles,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// Helper methods
func (u *UserEntity) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *UserEntity) HasPermission(resource string, action string) bool {
	permKey := fmt.Sprintf("%s:%s", resource, action)
	permissions := u.flattenPermissions()

	_, ok := permissions[permKey]
	return ok
}

func (u *UserEntity) flattenPermissions() map[string]struct{} {
	permMap := make(map[string]struct{})

	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			key := fmt.Sprintf("%s:%s", perm.Resource, perm.Action)
			permMap[key] = struct{}{}
		}
	}

	return permMap
}
