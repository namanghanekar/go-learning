package service

import (
	"fmt"

	"order-processing-system/notification-service/internal/email"
	"order-processing-system/notification-service/internal/sms"
)

type NotificationService struct {
	EmailSender *email.EmailSender
	SMSSender   *sms.SMSSender
}

func NewNotificationService(
	emailSender *email.EmailSender,
	smsSender *sms.SMSSender,
) *NotificationService {

	return &NotificationService{
		EmailSender: emailSender,
		SMSSender:   smsSender,
	}
}

func (s *NotificationService) SendNotification(
	orderID int,
) {

	message := fmt.Sprintf(
		"Your order %d has been confirmed",
		orderID,
	)

	s.EmailSender.SendEmail(
		"user@gmail.com",
		message,
	)

	s.SMSSender.SendSMS(
		"9999999999",
		message,
	)
}
