package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/responses"
)

// type key string

// const tokenUID key = "tokenuid"

// ExtractUID verifies idToken and extracts usertoken
func ExtractUID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("token")
		token, err := config.Client.VerifyIDToken(context.TODO(), idToken)

		if err != nil {
			log.Println("ExtractUID: ", err)
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "tokenuid", token)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
