package utils

import (
	"context"
	"log"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/mailgun/mailgun-go/v4"
)


func SendMail(link, name, email, mode string) error {

	var subject string

	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAINGUN_APIKEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := os.Getenv("MAILGUN_SENDER")
	if mode == "emailverification" {
		subject = "Verify your zode account"
	}else {
		subject = "Zode account password reset"
	}
	body := link
	recipient := email

	message := mg.NewMessage(sender, subject, body, recipient)

	if mode == "emailverification" {
		message.SetTemplate("emailverification")
		message.AddTemplateVariable("emailVerificationLink", link)
	}else {
		message.SetTemplate("passwordreset")
		message.AddTemplateVariable("passwordResetLink", link)
	}
	message.AddTemplateVariable("name", name)
	message.AddTemplateVariable("email", email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Error at SendPasswordResetLink sendmail.go : %v", err)
		return err
	}

	return nil
}

// SendEmailVerificationLink send link to verify the user's email
func SendEmailVerificationLink(name, email string) error {

	actionCodeSettings := auth.ActionCodeSettings{
		URL: "https://zod-frontend.herokuapp.com/login",
	}

	link, err := config.Client.EmailVerificationLinkWithSettings(context.Background(), email, &actionCodeSettings)
	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	err = SendMail(link, name, email, "emailverification")

	if err != nil {
		log.Printf("Error at SendEmailVerification sendMail.go : %v", err)
		return err
	}

	return nil
}

// SendPasswordResetLink send link to reset user's password
func SendPasswordResetLink(name, email string) error {
	actionCodeSettings := auth.ActionCodeSettings{
		URL: "https://zod-frontend.herokuapp.com/login",
	}

	link, err := config.Client.PasswordResetLinkWithSettings(context.Background(), email, &actionCodeSettings)

	if err != nil {
		log.Printf("Error at SendPasswordResetLink sendmail.go : %v", err)
		return err
	}

	err = SendMail(link, name, email, "passwordreset")

	if err != nil {
		log.Printf("Error at SendPasswordResetLink sendmail.go : %v", err)
		return err
	}

	return nil
}
