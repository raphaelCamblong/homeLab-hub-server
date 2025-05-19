package usecase

import (
	"errors"
	"strconv"

	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type UserUseCase interface {
	// Authentication
	Login(username, password string, ipAddress, userAgent string) (*entities.AuthResponse, error)
	Logout(userID string) error
	Register(username, password string) (*entities.AuthResponse, error)
	GetMe(userID string) (*entities.UserResponse, error)

	// User Management
	GetUser(id string) (*entities.UserResponse, error)
	ListUsers(page, pageSize int) ([]entities.UserResponse, int64, error)
	UpdateUserStatus(id string, active bool) error

	// Role Management
	GetUserRoles(userID string) ([]entities.Role, error)
	AssignRole(userID string, roleID string) error
	RemoveRole(userID string, roleID string) error
	ListRoles() ([]entities.Role, error)

	// Permission Management
	ListPermissions() ([]entities.Permission, error)
	GetUserPermissions(userID string) ([]entities.Permission, error)
	AssignPermission(roleID string, permissionID string) error
	RemovePermission(roleID string, permissionID string) error

	// API Keys
	CreateAPIKey(userID string, name string) (*entities.APIKey, error)
	ListAPIKeys(userID string) ([]entities.APIKey, error)
	RevokeAPIKey(userID string, keyID string) error

	// Audit
	GetAuditLogs(userID string, page, pageSize int) ([]entities.AuditLog, int64, error)
}

type userUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (u *userUseCase) Login(username, password string, ipAddress, userAgent string) (*entities.AuthResponse, error) {
	cred := entities.UserCredentials{
		Username: username,
		Password: password,
	}
	return u.userRepo.Login(cred, ipAddress, userAgent)
}

func (u *userUseCase) Logout(userID string) error {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return u.userRepo.Logout(uint(id))
}

func (u *userUseCase) Register(username, password string) (*entities.AuthResponse, error) {
	cred := entities.UserCredentials{
		Username: username,
		Password: password,
	}
	return u.userRepo.Register(cred)
}

func (u *userUseCase) GetUser(id string) (*entities.UserResponse, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.GetUser(uint(uid))
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUseCase) ListUsers(page, pageSize int) ([]entities.UserResponse, int64, error) {
	return u.userRepo.ListUsers(page, pageSize)
}

func (u *userUseCase) UpdateUserStatus(id string, active bool) error {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return u.userRepo.SetUserActive(uint(uid), active)
}

func (u *userUseCase) GetUserRoles(userID string) ([]entities.Role, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.GetUser(uint(uid))
	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}

func (u *userUseCase) AssignRole(userID string, roleID string) error {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		return errors.New("invalid role ID")
	}

	return u.userRepo.AssignRole(uint(uid), uint(rid))
}

func (u *userUseCase) RemoveRole(userID string, roleID string) error {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		return errors.New("invalid role ID")
	}

	return u.userRepo.RemoveRole(uint(uid), uint(rid))
}

func (u *userUseCase) ListRoles() ([]entities.Role, error) {
	return u.userRepo.ListRoles()
}

func (u *userUseCase) ListPermissions() ([]entities.Permission, error) {
	return u.userRepo.ListPermissions()
}

func (u *userUseCase) GetUserPermissions(userID string) ([]entities.Permission, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.GetUserPermissions(uint(uid))
}

func (u *userUseCase) AssignPermission(roleID string, permissionID string) error {
	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		return errors.New("invalid role ID")
	}

	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		return errors.New("invalid permission ID")
	}

	return u.userRepo.AssignPermission(uint(rid), uint(pid))
}

func (u *userUseCase) RemovePermission(roleID string, permissionID string) error {
	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		return errors.New("invalid role ID")
	}

	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		return errors.New("invalid permission ID")
	}

	return u.userRepo.RemovePermission(uint(rid), uint(pid))
}

func (u *userUseCase) CreateAPIKey(userID string, name string) (*entities.APIKey, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.CreateAPIKey(uint(uid), name, nil, nil)
}

func (u *userUseCase) ListAPIKeys(userID string) ([]entities.APIKey, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.ListAPIKeys(uint(uid))
}

func (u *userUseCase) RevokeAPIKey(userID string, keyID string) error {
	kid, err := strconv.ParseUint(keyID, 10, 32)
	if err != nil {
		return errors.New("invalid key ID")
	}

	return u.userRepo.RevokeAPIKey(uint(kid))
}

func (u *userUseCase) GetAuditLogs(userID string, page, pageSize int) ([]entities.AuditLog, int64, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	return u.userRepo.ViewAuditLogs(uint(uid), page, pageSize)
}

func (u *userUseCase) GetMe(userID string) (*entities.UserResponse, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.GetMe(uint(id))
}
