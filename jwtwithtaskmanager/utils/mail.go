package utils

import (
	"fmt"
	"net/smtp"
)

func SendOTPEmail(to string, otp string) {

	from := "fitnessprojectnnnp25@gmail.com"
	password := "nlug crej zjuc vbdj"
	msg := []byte("Subject: OTP Verification\n\nYour OTP is: " + otp)
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		msg,
	)
	if err != nil {
		fmt.Println("Email error:", err)
		return
	}
	fmt.Println("OTP sent to:", to)
}
