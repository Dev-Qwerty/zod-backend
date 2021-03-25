package utilities

import (
	"context"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

func SendEmailVerificationLink(email string) error {
	link, err := config.Client.EmailVerificationLink(context.Background(), email)
	if err != nil {
		return err
	}

	err = config.SendMail(link, email)
	if err != nil {
		return err
	}

	return nil
}
