package sms

import "log"

type SMSSender struct {
}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

func (s *SMSSender) SendSMS(
	phone string,
	message string,
) {

	log.Printf(
		"Sending SMS to %s : %s\n",
		phone,
		message,
	)
}
