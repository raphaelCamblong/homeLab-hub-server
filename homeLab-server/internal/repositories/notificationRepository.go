package repositories

import (
	"context"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/streaming"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
)

type NotificationRepository interface {
	List() ([]entities.Notification, error)
	Get(id string) (entities.Notification, error)
	Create(notification entities.Notification) (entities.Notification, error)
	Delete(id string) error
	Subscribe(ctx context.Context, clientID string) (<-chan entities.Notification, error)
}

type notificationRepository struct {
	db        database.Database
	streamHub *streaming.StreamHub
}

func NewNotificationRepository(db database.Database, streamHub *streaming.StreamHub) NotificationRepository {
	return &notificationRepository{
		db:        db,
		streamHub: streamHub,
	}
}

func (r *notificationRepository) List() ([]entities.Notification, error) {
	var notifications []entities.Notification
	result := r.db.Get().Order("created_at DESC").Find(&notifications)
	return notifications, result.Error
}

func (r *notificationRepository) Get(id string) (entities.Notification, error) {
	var notification entities.Notification
	result := r.db.Get().First(&notification, "id = ?", id)
	return notification, result.Error
}

func (r *notificationRepository) Create(notification entities.Notification) (entities.Notification, error) {
	result := r.db.Get().Create(&notification)
	if result.Error != nil {
		return entities.Notification{}, result.Error
	}
	return notification, nil
}

func (r *notificationRepository) Delete(id string) error {
	result := r.db.Get().Delete(&entities.Notification{}, "id = ?", id)
	return result.Error
}

func (r *notificationRepository) Subscribe(ctx context.Context, clientID string) (<-chan entities.Notification, error) {
	eventCh, err := r.streamHub.Subscribe(ctx, streaming.NotificationEvent)
	if err != nil {
		return nil, err
	}

	notificationCh := make(chan entities.Notification)

	go func() {
		defer close(notificationCh)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventCh:
				if !ok {
					return
				}
				if notification, ok := event.Data.(entities.Notification); ok {
					notificationCh <- notification
					if _, err := r.Create(notification); err != nil {
						logrus.Info("Failed to create Notification in DB")
					}
				}
			}
		}
	}()

	return notificationCh, nil
}
