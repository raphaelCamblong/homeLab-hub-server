package entities

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:255" json:"username"`
	Password string `json:"-"`
}

type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtToken string

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
