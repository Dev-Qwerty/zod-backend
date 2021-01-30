package controllers

import (
	"encoding/json"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/models"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/responses"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	project.Teamlead = userDetails.DisplayName

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

// GetProjectsHandler is the handler for /getprojects route
func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	projects, err := project.GetProjects(tokenStruct.Claims["email"].(string))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	responses.JSON(w, http.StatusOK, projects)
}

// AddProjectMembersHandler creates handler for /addnewprojectmembers route
func AddProjectMembersHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = project.AddProjectMembers(tokenStruct.Claims["email"].(string))
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	responses.JSON(w, http.StatusOK, nil)
}

// AcceptInviteHandler creates handler for /acceptinvites route
func AcceptInviteHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	userDetails, _ := utils.GetUserDetails(tokenStruct.UID) // todo : handle error

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = project.AcceptInvite(userDetails)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, nil)
}

// LeaveProjectHandler creates handler for /leaveproject route
func LeaveProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	project := &models.Project{}
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = project.LeaveProject(tokenStruct.Claims["email"].(string))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, nil)
}

func RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	type Details struct {
		ProjectID primitive.ObjectID `json:"projectID,omitempty"`
		Email     string             `json:"email,omitempty"`
	}
	var detail *Details
	err := json.NewDecoder(r.Body).Decode(&detail)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = models.RemoveProjectMember(tokenStruct.Claims["email"].(string), detail.Email, detail.ProjectID)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, nil)
}
