package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type NasHandler struct {
	nasUseCase usecase.NasUseCase
}

func NewNasHandler(nasUseCase usecase.NasUseCase) *NasHandler {
	return &NasHandler{nasUseCase: nasUseCase}
}

func (h *NasHandler) ListStoragePools(ctx *gin.Context) {
	pools, err := h.nasUseCase.ListStoragePools()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pools)
}

func (h *NasHandler) GetStoragePoolDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	details, err := h.nasUseCase.GetStoragePoolDetails(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, details)
}

func (h *NasHandler) GetZFSMetrics(ctx *gin.Context) {
	metrics, err := h.nasUseCase.GetZFSMetrics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}

func (h *NasHandler) BrowseFilesystem(ctx *gin.Context) {
	files, err := h.nasUseCase.BrowseFilesystem()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, files)
}

func (h *NasHandler) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.nasUseCase.UploadFile(file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (h *NasHandler) DownloadFile(ctx *gin.Context) {
	fileID := ctx.Query("id")
	file, err := h.nasUseCase.DownloadFile(fileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, file)
}

func (h *NasHandler) TriggerBackupJob(ctx *gin.Context) {
	if err := h.nasUseCase.TriggerBackupJob(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Backup job triggered successfully"})
} 