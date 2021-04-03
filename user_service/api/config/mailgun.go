package config

import (
	"context"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendMail(link string, email string) error {

	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAINGUN_APIKEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := os.Getenv("MAILGUN_SENDER")
	subject := "Verify your zode account"
	body := link
	recipient := email

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	return nil
}
