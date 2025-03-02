package usecase

import (
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type UserUseCase interface {
	Login(credentials Credentials) (string, error)
	Logout() error
	ViewUserRoles() ([]Role, error)
	ViewAuditLogs() ([]AuditLog, error)
}

type userUseCase struct {
	userRepository repositories.UserRepository
}

func NewUserUseCase(userRepository repositories.UserRepository) UserUseCase {
	return &userUseCase{userRepository: userRepository}
}

func (u *userUseCase) Login(credentials Credentials) (string, error) {
	return u.userRepository.Login(credentials)
}

func (u *userUseCase) Logout() error {
	return u.userRepository.Logout()
}

func (u *userUseCase) ViewUserRoles() ([]Role, error) {
	return u.userRepository.ViewUserRoles()
}

func (u *userUseCase) ViewAuditLogs() ([]AuditLog, error) {
	return u.userRepository.ViewAuditLogs()
} 