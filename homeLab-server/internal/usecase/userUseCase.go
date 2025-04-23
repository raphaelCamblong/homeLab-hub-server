package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
)

type UserUseCase interface {
	Login(username, password string) (User, error)
	Logout(userID string) error
	GetPermissions(userID string) ([]Permission, error)
	ViewUserRoles() ([]Role, error)
	ViewAuditLogs() ([]AuditLog, error)
	Register(credentials authentication.UserCredentials) (*authentication.JwtToken, error)
}

type userUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (u *userUseCase) Login(username, password string) (User, error) {
	cred := authentication.UserCredentials{Username: username, Password: password}
	token, err := u.userRepo.Login(cred)
	if err != nil {
		return User{}, err
	}
	return User{Token: token}, nil
}

func (u *userUseCase) Logout(userID string) error {
	return u.userRepo.Logout(userID)
}

func (u *userUseCase) GetPermissions(userID string) ([]Permission, error) {
	return []Permission{}, nil
}

func (u *userUseCase) ViewUserRoles() ([]Role, error) {
	return u.userRepo.ViewUserRoles()
}

func (u *userUseCase) ViewAuditLogs() ([]AuditLog, error) {
	return u.userRepo.ViewAuditLogs()
}

func (u *userUseCase) Register(credentials authentication.UserCredentials) (*authentication.JwtToken, error) {
	return u.userRepo.Register(credentials)
}
