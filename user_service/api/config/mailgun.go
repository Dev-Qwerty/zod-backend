package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendMail(link ,name, email, subjectline string) error {

	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAINGUN_APIKEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := os.Getenv("MAILGUN_SENDER")
	subject := subjectline
	body := link
	recipient := email

	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("emailverification")
	message.AddTemplateVariable("name", name)
	message.AddTemplateVariable("emailVerificationLink", link)
	message.AddTemplateVariable("email", email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Error at SendMail mailgun.go : %v", err)
		return err
	}

	return nil
}
