package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/responses"
)

// Add projects of user to db
func NewProject(w http.ResponseWriter, r *http.Request) {

	// var project map[string]interface{} //Message queues
	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Printf("Error at Decode user project_controller.go: %v", err)
	}

	/* Message Queues
	projectByte, err := json.Marshal(project)
	if err != nil {
		log.Printf("Error at NewProject project_controller.go: %v", err)
	}
	*/

	// err = messageQueues.Produce("Create project", projectByte)
	err = models.AddProject(project)
	if err != nil {
		fmt.Printf("Failed to save project: %v", err)
		responses.ERROR(w, http.StatusInternalServerError, err)
	} else {
		responses.JSON(w, http.StatusOK, nil)
	}

}

// Update project role of users
func UpdateProject(w http.ResponseWriter, r *http.Request) {

	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Printf("Failed decoding project: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		err = models.UpdateProject(project)

		if err != nil {
			fmt.Printf("Failed to update project: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		} else {
			responses.JSON(w, http.StatusOK, nil)
		}
	}
}

// DeleteProject delete projects associated with a user
func DeleteProject(w http.ResponseWriter, r *http.Request) {

	project := models.Project{}

	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		log.Printf("Failed decoding project: %v", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	} else {
		err := models.DeleteProject(project)

		if err != nil {
			fmt.Printf("Failed to delete project: %v", err)
			responses.ERROR(w, http.StatusInternalServerError, err)
		} else {
			responses.JSON(w, http.StatusNoContent, nil)
		}
	}
}
