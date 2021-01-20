package controllers

import (
	"encoding/json"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/models"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/responses"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/utils"
)

// CreateProjectHandler creates the handler for createproject route
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	userDetails, err := utils.GetUserDetails(tokenStruct.UID)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	project := &models.Project{}
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// creating member model of the user created the project with an owner role
	member := &models.Member{}
	member.Name = userDetails.DisplayName
	member.Email = userDetails.Email
	member.UserID = userDetails.UID
	member.Role = "Owner"

	project.Members = &[]models.Member{}

	*project.Members = append(*project.Members, *member)

	projectID, err := project.CreateProject()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return

	}

	responses.JSON(w, http.StatusOK, projectID)
}
