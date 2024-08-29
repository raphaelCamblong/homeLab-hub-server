package handlers

import (
	"github.com/gin-gonic/gin"
	entities "homelab.com/homelab-server/homeLab-server/internal/entities/authentication"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
	"net/http"
)

type AuthenticationHandler struct {
	authenticationUseCase usecase.AuthenticationUseCase
}

func NewAuthenticationHandler(authenticationUseCase usecase.AuthenticationUseCase) *AuthenticationHandler {
	return &AuthenticationHandler{authenticationUseCase: authenticationUseCase}
}

func (h *AuthenticationHandler) Register(c *gin.Context) {
	var credentials entities.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authenticationUseCase.Register(credentials)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"x-auth-token": token})
}

func (h *AuthenticationHandler) Login(c *gin.Context) {
	var credentials entities.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authenticationUseCase.Login(credentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"x-auth-token": token})
}
