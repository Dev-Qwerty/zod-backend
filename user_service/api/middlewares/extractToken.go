package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/responses"
)

// ExtractToken extracts the use from the idToken
func ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		u, err := config.Client.VerifyIDToken(context.TODO(), token)
		if err != nil {
			log.Printf("Error at ExtractToken extractToken.go : %v", err)
			responses.ERROR(w, http.StatusUnauthorized, err)
		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "uid", u.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
