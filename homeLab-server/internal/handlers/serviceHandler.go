package handlers

import (
	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
	"net/http"
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

type UserRequest struct {
	Id string `uri:"id"`
}

func (s *ServiceHandler) GetServiceById(ctx *gin.Context) {
	var u UserRequest

	if err := ctx.ShouldBindUri(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.serviceUseCase.GetServiceById(u.Id)
	if err != nil {
		if err.Error() == "Service not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	ctx.JSON(http.StatusOK, data)
}
