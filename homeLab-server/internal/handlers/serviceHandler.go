package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type ServiceHandler struct {
	serviceUseCase usecase.ServiceUseCase
}

func NewServiceHandler(serviceUseCase usecase.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{serviceUseCase: serviceUseCase}
}

// ListServices godoc
// @Summary List all services
// @Description Get all services
// @Tags services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.ServiceEntity
// @Router /api/v1/services [get]
func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.serviceUseCase.ListServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetService godoc
// @Summary Get a service by ID
// @Description Get a single service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Security BearerAuth
// @Success 200 {object} entities.ServiceEntity
// @Router /api/v1/services/{id} [get]
func (h *ServiceHandler) GetService(c *gin.Context) {
	id := c.Param("id")
	service, err := h.serviceUseCase.GetService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service with the provided information
// @Tags services
// @Accept json
// @Produce json
// @Param service body entities.ServiceEntity true "Service object"
// @Security BearerAuth
// @Success 201 {object} entities.ServiceEntity
// @Router /api/v1/services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service entities.ServiceEntity
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	createdService, err := h.serviceUseCase.CreateService(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, createdService)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update a service with the provided information
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body entities.ServiceEntity true "Service object"
// @Security BearerAuth
// @Success 200 {object} entities.ServiceEntity
// @Router /api/v1/services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")
	var service entities.ServiceEntity
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedService, err := h.serviceUseCase.UpdateService(id, service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, updatedService)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Security BearerAuth
// @Success 200 {object} string "Service deleted successfully"
// @Router /api/v1/services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")
	if err := h.serviceUseCase.DeleteService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	// @Security BearerAuth
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
