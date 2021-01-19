package middlewares

import (
	"context"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/config"
)

type key string

const tokenKey key = "token"

// ExtractToken verifies idToken and extracts usertoken
func ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _ := config.Client.VerifyIDToken(context.TODO(), "idToken")
		ctx := r.Context()
		ctx = context.WithValue(ctx, tokenKey, token)

		next.ServeHTTP(w, r)
	})
}
