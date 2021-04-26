package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/responses"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/utils"
)

// SignUp creates new user
func SignUp(w http.ResponseWriter, r *http.Request) {

	user := models.FirebaseUser{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		err := models.CreateNewUser(user)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		} else {
			responses.JSON(w, http.StatusCreated, nil)
		}
	}
}

// Update updates the user data
func Update(w http.ResponseWriter, r *http.Request) {

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)

	user := models.UpdatedUser{}
	user.ID = r.Context().Value("uid").(string)
	
	if v, found := data["email"]; found {
		user.Email = v
	}

	if v, found := data["fname"]; found {
		user.Fname = v
	}

	if v, found := data["lname"]; found {
		user.Lname = v
	}
	
	if err != nil {
		log.Printf("Failed decoding user: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		err := models.UpdateUser(user)
		if err != nil {
			fmt.Printf("Failed to update user: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		} else {
			responses.JSON(w, http.StatusOK, nil)
		}
	}
}

// Delete removes the user
func Delete(w http.ResponseWriter, r *http.Request) {

	uid := r.Context().Value("uid").(string)

	err := models.DeleteUser(uid)
	if err != nil {
		fmt.Printf("Failed to delete user: %v", err)
		responses.ERROR(w, http.StatusInternalServerError, err)
	} else {
		responses.JSON(w, http.StatusNoContent, nil)
	}
}

// ResendEmail sends the verification email
func ResendEmail(w http.ResponseWriter, r *http.Request) {
	var email map[string]string

	err := json.NewDecoder(r.Body).Decode(&email)

	if err != nil {
		fmt.Printf("Failed to send email : %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		u, err := config.Client.GetUserByEmail(context.Background(), email["email"])
		if err != nil {
			log.Printf("Error at ResendEmail user_controller.go: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		err = utils.SendEmailVerificationLink(u.UserInfo.DisplayName, email["email"])
		if err != nil {
			log.Printf("Error at ResendEmail user_controller.go : %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		responses.JSON(w, http.StatusOK, nil)
	}
}

// ResetPassword send a link to reset the password
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var email map[string]string

	err := json.NewDecoder(r.Body).Decode(&email)

	if err != nil {
		fmt.Printf("Failed to send email: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		u, err := config.Client.GetUserByEmail(context.Background(), email["email"])
		if err != nil {
			log.Printf("Error at ResetPassword user_controller.go: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		err = utils.SendPasswordResetLink(u.UserInfo.DisplayName, email["email"])

		if err != nil {
			log.Printf("Error at ResetPassword user_controller.go : %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		responses.JSON(w, http.StatusOK, nil)
	}
}
