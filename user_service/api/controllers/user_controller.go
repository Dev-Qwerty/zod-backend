package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
)

// UserSignUp creates new user
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	user := models.FirebaseUser{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %v", err)
	} else {
		models.CreateNewUser(user)
	}
}
