package middlewares

import (
	"context"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
)

// ExtractToken extracts the use from the idToken
func ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		u, err := config.Client.VerifyIDToken(context.TODO(), token)
		if err != nil {

		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "token", u.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
