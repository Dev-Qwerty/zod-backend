package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	data := models.UpdatedUser{}
	data.ID = r.Context().Value("uid").(string)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Failed decoding user: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		err := models.UpdateUser(data)
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
		err = utils.SendEmailVerificationLink(email["email"])
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
		err = utils.SendPasswordResetLink(email["email"])

		if err != nil {
			log.Printf("Error at ResendEmail user_controller.go : %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		responses.JSON(w, http.StatusOK, nil)
	}
}
