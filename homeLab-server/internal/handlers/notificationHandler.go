package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type NotificationHandler struct {
	notificationUseCase usecase.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

// ListNotifications godoc
// @Summary List all notifications
// @Description Get all notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Notification
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	notifications, err := h.notificationUseCase.ListNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// StreamNotifications godoc
// @Summary Stream notifications using SSE
// @Description Establishes a Server-Sent Events connection to stream notifications
// @Tags notifications
// @Accept json
// @Produce text/event-stream
// @Security BearerAuth
// @Success 200 {string} string "SSE connection established"
// @Router /api/v1/notifications/stream [get]
func (h *NotificationHandler) StreamNotifications(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	clientID := uuid.New().String()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	updates, err := h.notificationUseCase.Subscribe(ctx, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe to job updates"})
		return
	}

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case notification, ok := <-updates:
			if !ok {
				return false
			}
			c.SSEvent("notification", notification)
			return true
		case <-time.After(30 * time.Second):
			c.SSEvent("heartbeat", map[string]string{
				"timestamp": time.Now().Format(time.RFC3339),
			})
			c.Writer.Flush()
			return true
		}
	})
}
