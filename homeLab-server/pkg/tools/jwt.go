package tools

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

func GenerateJWT(username string, id uint, duration int, secret string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(duration) * time.Minute)
	secretKey := []byte(secret)

	claims := &entities.Claims{
		Username: username,
		UserID:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseJWT(tokenStr string, secret string) (*entities.Claims, error) {
	claims := &entities.Claims{}
	secretKey := []byte(secret)

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
