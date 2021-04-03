package utilities

import (
	"context"
	"log"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

func SendEmailVerificationLink(email string) error {
	link, err := config.Client.EmailVerificationLink(context.Background(), email)
	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	err = config.SendMail(link, email)
	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	return nil
}
