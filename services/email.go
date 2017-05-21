package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/utils"
)

func SendResetTokenEmail(to string, token string) error {
	from := "sofyan.h.ahmad@gmail.com"
	pass := os.Getenv("EMAILPASSWORD")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Token\n\n" +
		"You requested to reset yoour password, ignore if you never \n" +
		"to reset your password follow this link: " + utils.BaseUrl + "user/changepassword/do?email=" + to + "&t=" + token

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print(fmt.Sprintf("Reset token sent to %s : token %s", to, token))
	return err
}
