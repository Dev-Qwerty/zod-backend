package middlewares

import (
	"context"
	"log"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

// ExtractToken verifies the token and returns UID
func ExtractToken() string {
	token, err := config.Client.VerifyIDToken(context.Background(), "idToken")
	if err != nil {
		log.Printf("Failed to verify token: %v", err)
	}
	return token.UID
}
