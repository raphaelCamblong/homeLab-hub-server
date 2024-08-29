package tools

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
	"time"
)

func GenerateJWT(username string) (string, error) {
	//TODO: Change expiration time to a more secure value
	expirationTime := time.Now().Add(3000 * time.Hour)
	//TODO: Change secret key to a more secure value
	secretKey := []byte("secret")

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
	//TODO: Change secret key to a more secure value
	secretKey := []byte("secret")

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
