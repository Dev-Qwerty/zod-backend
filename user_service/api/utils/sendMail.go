package utils

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

// SendEmailVerificationLink send link to verify the user's email
func SendEmailVerificationLink(email string) error {

	actionCodeSettings := auth.ActionCodeSettings{
		URL: "http://localhost:3000/login",
	}

	link, err := config.Client.EmailVerificationLinkWithSettings(context.Background(), email, &actionCodeSettings)
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
	actionCodeSettings := auth.ActionCodeSettings{
		URL: "http://localhost:3000/login",
	}

	link, err := config.Client.PasswordResetLinkWithSettings(context.Background(), email, &actionCodeSettings)

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
