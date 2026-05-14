package repository

import (
	"order-processing-system/notification-service/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(
	db *gorm.DB,
) *NotificationRepository {

	return &NotificationRepository{
		DB: db,
	}
}

func (r *NotificationRepository) Create(
	notification *model.Notification,
) error {

	return r.DB.Create(notification).Error
}
