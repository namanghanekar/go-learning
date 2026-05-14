package service

import (
	"fmt"

	"order-processing-system/notification-service/internal/model"
	"order-processing-system/notification-service/internal/repository"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func NewNotificationService(
	repo *repository.NotificationRepository,
) *NotificationService {

	return &NotificationService{
		Repo: repo,
	}
}

func (s *NotificationService) SendNotification(
	userID int,
	message string,
) error {

	fmt.Println(
		"Notification sent",
	)

	notification := &model.Notification{
		UserID:  userID,
		Message: message,
		Type:    "ORDER",
	}

	return s.Repo.Create(
		notification,
	)
}
