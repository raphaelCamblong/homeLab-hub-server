package usecase

import (
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type AuthenticationUseCase interface {
	Register(entities.UserCredentials) (*entities.JwtToken, error)
	Login(entities.UserCredentials) (*entities.JwtToken, error)
}

type authenticationUseCase struct {
	repository repositories.AuthenticationRepository
}

func NewAuthenticationUseCase(repository repositories.AuthenticationRepository) AuthenticationUseCase {
	return &authenticationUseCase{repository: repository}
}

func (a *authenticationUseCase) Register(cred entities.UserCredentials) (*entities.JwtToken, error) {
	return a.repository.Register(cred)
}

func (a *authenticationUseCase) Login(cred entities.UserCredentials) (*entities.JwtToken, error) {
	return a.repository.Login(cred)
}
