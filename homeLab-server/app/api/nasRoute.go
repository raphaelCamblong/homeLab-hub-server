package api

import (
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/internal/handlers"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

func NasRoutes(infra *infrastructure.Infrastructure, repo *Repositories) error {
	router := infra.Router.Get().Group("/api/v1/nas")

	handler := handlers.NewNasHandler(usecase.NewNasUseCase(repo.Nas))

	router.GET("/storage", handler.ListStoragePools)
	router.GET("/storage/{id}", handler.GetStoragePoolDetails)
	router.GET("/zfs", handler.GetZFSMetrics)
	router.GET("/files", handler.BrowseFilesystem)
	router.POST("/files/upload", handler.UploadFile)
	router.GET("/files/download", handler.DownloadFile)
	router.POST("/backup", handler.TriggerBackupJob)

	return nil
}
