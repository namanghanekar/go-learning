package utils

import (
	"fmt"
	"net/smtp"
)

func SendOTPEmail(to string, otp string) error {

	from := "your-email@gmail.com"
	password := "your-app-password"

	msg := []byte("Subject: OTP Verification\n\nYour OTP is: " + otp)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		msg,
	)

	if err != nil {
		return err
	}

	fmt.Println("✅ OTP sent to", to)
	return nil
}
