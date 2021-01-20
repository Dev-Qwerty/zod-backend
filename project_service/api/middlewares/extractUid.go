package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/config"
)

// type key string

// const tokenUID key = "tokenuid"

// ExtractUID verifies idToken and extracts usertoken
func ExtractUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("token")
		token, err := config.Client.VerifyIDToken(context.TODO(), idToken)
		if err != nil {
			fmt.Println(err)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "tokenuid", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
