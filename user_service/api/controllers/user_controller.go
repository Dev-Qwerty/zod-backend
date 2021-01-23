package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
)

// SignUp creates new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	user := models.FirebaseUser{}

	user.Email = "email"

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %v", err)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
	} else {
		err := models.CreateNewUser(user)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

// Update updates the user data
func Update(w http.ResponseWriter, r *http.Request) {

	data := models.UpdatedUser{}
	ctx := r.Context()
	data.Id = ctx.Value("token").(string)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Failed decoding user: %v", err)
		http.Error(w, "Failed to update user", http.StatusBadRequest)
	} else {
		err := models.UpdateUser(data)
		if err != nil {
			fmt.Printf("Failed to update user: %v", err)
			http.Error(w, "Failed to update user", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
