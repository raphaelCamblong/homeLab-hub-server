package tools

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"homelab.com/homelab-server/homeLab-server/app/config"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
	"time"
)

func GenerateJWT(username string) (string, error) {
	cfg := config.GetConfig()
	expirationTime := time.Now().Add(time.Duration(cfg.App.Security.JWT.Expiration) * time.Minute)
	secretKey := []byte(cfg.App.Security.JWT.Secret)

	claims := &entities.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseJWT(tokenStr string) (*entities.Claims, error) {
	claims := &entities.Claims{}
	secretKey := []byte(config.GetConfig().App.Security.JWT.Secret)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
