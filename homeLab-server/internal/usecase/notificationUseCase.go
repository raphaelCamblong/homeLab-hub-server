package usecase

import (
	"context"

	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type NotificationUseCase interface {
	ListNotifications() ([]entities.Notification, error)
	GetNotification(id string) (entities.Notification, error)
	CreateNotification(notification entities.Notification) (entities.Notification, error)
	DeleteNotification(id string) error
	Subscribe(ctx context.Context, clientID string) (<-chan entities.Notification, error)
}

type notificationUseCase struct {
	notificationRepo repositories.NotificationRepository
}

func NewNotificationUseCase(notificationRepo repositories.NotificationRepository) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
	}
}

func (u *notificationUseCase) ListNotifications() ([]entities.Notification, error) {
	return u.notificationRepo.List()
}

func (u *notificationUseCase) GetNotification(id string) (entities.Notification, error) {
	return u.notificationRepo.Get(id)
}

func (u *notificationUseCase) CreateNotification(notification entities.Notification) (entities.Notification, error) {
	return u.notificationRepo.Create(notification)
}

func (u *notificationUseCase) DeleteNotification(id string) error {
	return u.notificationRepo.Delete(id)
}

func (u *notificationUseCase) Subscribe(ctx context.Context, clientID string) (<-chan entities.Notification, error) {
	return u.notificationRepo.Subscribe(ctx, clientID)
}
