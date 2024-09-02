package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type ServiceHandler struct {
	serviceUseCase usecase.ServiceUseCase
}

func NewServiceHandler(serviceUseCase usecase.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{serviceUseCase: serviceUseCase}
}

func (s *ServiceHandler) GetAllService(ctx *gin.Context) {
	data, err := s.serviceUseCase.GetAllService()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (s *ServiceHandler) GetServiceById(ctx *gin.Context) {
	_, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	data, err := s.serviceUseCase.GetServiceById(ctx.Param("id"))
	if err != nil {
		if err.Error() == "Service not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	ctx.JSON(http.StatusOK, data)
}
