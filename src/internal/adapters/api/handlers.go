package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourproject/internal/domain"
	"github.com/yourproject/internal/usecase"
)

type Handler struct {
	useCase usecase.ServiceUsecase
}

func NewHandler(uc usecase.ServiceUsecase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) CreateService(c *gin.Context) {
	var service domain.Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdService, err := h.useCase.CreateService(c.Request.Context(), service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdService)
}
