package repositories

import (
	"errors"
	"gorm.io/gorm"
	"homelab.com/homelab-server/homeLab-server/cmd/tools"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
)

type AuthenticationRepository interface {
	Register(entities.UserCredentials) (*entities.JwtToken, error)
	Login(entities.UserCredentials) (*entities.JwtToken, error)
}

type authenticationRepository struct {
	db database.Database
}

func NewAuthenticationRepository(db database.Database) AuthenticationRepository {
	return &authenticationRepository{db: db}
}

func (a *authenticationRepository) Register(cred entities.UserCredentials) (*entities.JwtToken, error) {
	var existingUser entities.UserEntity
	if err := a.db.GetDb().Where("username = ?", cred.Username).First(&existingUser).Error; err == nil {
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

	if err := a.db.GetDb().Create(&newUser).Error; err != nil {
		return nil, err
	}

	return a.Login(cred)
}

func (a *authenticationRepository) Login(cred entities.UserCredentials) (*entities.JwtToken, error) {
	var user entities.UserEntity
	if err := a.db.GetDb().Where("username = ?", cred.Username).First(&user).Error; err != nil {
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
