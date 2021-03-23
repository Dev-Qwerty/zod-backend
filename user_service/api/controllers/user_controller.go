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
	data.ID = r.Context().Value("uid").(string)

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

// Delete removes the user
func Delete(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("uid").(string)

	err := models.DeleteUser(uid)
	if err != nil {
		fmt.Printf("Failed to delete user: %v", err)
		http.Error(w, "Failed to delete user", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Add projects of user to db
func NewProject(w http.ResponseWriter, r *http.Request) {
	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Printf("Failed decoding project: %v", err)
		http.Error(w, "Failed to save project", http.StatusBadRequest)
	} else {
		err = models.AddProject(project)
		if err != nil {
			fmt.Printf("Failed to save project: %v", err)
			http.Error(w, "Failed to save project", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// Update project role of users
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Printf("Failed decoding project: %v", err)
		http.Error(w, "Failed to update project", http.StatusBadRequest)
	} else {
		err = models.UpdateProject(project)

		if err != nil {
			fmt.Printf("Failed to update project: %v", err)
			http.Error(w, "Failed to update project", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		log.Printf("Failed decoding project: %v", err)
		http.Error(w, "Failed to delete project", http.StatusBadRequest)
	} else {
		err := models.DeleteProject(project)

		if err != nil {
			fmt.Printf("Failed to delete project: %v", err)
			http.Error(w, "Failed to delete project", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
