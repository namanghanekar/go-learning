package email

import "log"

type EmailSender struct {
}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (e *EmailSender) SendEmail(
	to string,
	message string,
) {

	log.Printf(
		"Sending EMAIL to %s : %s\n",
		to,
		message,
	)
}
