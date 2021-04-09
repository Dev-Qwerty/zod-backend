package utils

import (
	"context"
	"log"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

// SendEmailVerificationLink send link to verify the user's email
func SendEmailVerificationLink(email string) error {
	link, err := config.Client.EmailVerificationLink(context.Background(), email)
	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	err = config.SendMail(link, email, "Verify your zode account")
	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	return nil
}

// SendPasswordResetLink send link to reset user's password
func SendPasswordResetLink(email string) error {
	link, err := config.Client.PasswordResetLink(context.Background(), email)

	if err != nil {
		log.Printf("Error at SendPasswordResetLink sendmail.go : %v", err)
		return err
	}

	err = config.SendMail(link, email, "Reset your zode account password")
	if err != nil {
		log.Printf("Error at SendPasswordResetLink sendmail.go : %v", err)
		return err
	}

	return nil
}
