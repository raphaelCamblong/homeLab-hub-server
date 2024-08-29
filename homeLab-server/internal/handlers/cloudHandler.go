package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type CloudHandler struct {
	clouduseCase usecase.CloudUseCase
}

func NewCloudHandler(clouduseCase usecase.CloudUseCase) *CloudHandler {
	return &CloudHandler{clouduseCase: clouduseCase}
}

func (h *CloudHandler) GetVmsData(ctx *gin.Context) {
	data, err := h.clouduseCase.GetVmsData()
	if err != nil {
		logrus.Errorf("failed to retrieve VMs data: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (h *CloudHandler) GetHostData(ctx *gin.Context) {
	data, err := h.clouduseCase.GetHostData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
