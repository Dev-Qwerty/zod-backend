package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
)

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
			http.Error(w, "Failed to save project", http.StatusInternalServerError)
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
			http.Error(w, "Failed to update project", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// DeleteProject delete projects associated with a user
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
			http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
