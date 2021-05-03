package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/models"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/responses"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/utils"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// CreateProjectHandler creates the handler for createproject route
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	userDetails, err := utils.GetUserDetails(tokenStruct.UID)

	if err != nil {
		log.Println("CreateProjectHandler: <Error while fetching user details> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	project := &models.Project{}
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Println("CreateProjectHandler: <Error while decoding request Body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	project.Teamlead = userDetails.DisplayName

	// creating member model of the user created the project with an Admin role
	member := &models.Member{}
	member.Name = userDetails.DisplayName
	member.Email = userDetails.Email
	member.UserID = userDetails.UID
	member.Role = "Admin"
	MemberID := uuid.NewV4().String()
	member.MemberID = MemberID[24:]

	project.Members = &[]models.Member{}

	*project.Members = append(*project.Members, *member)
	for i := 0; i < len(*project.PendingInvites); i++ {
		(*project.PendingInvites)[i].InvitedBy = userDetails.DisplayName
		(*project.PendingInvites)[i].Name, err = utils.GetUserDetailsByEmail((*project.PendingInvites)[i].Email)
		if err != nil {
			log.Println("CreateProjectHandler: <Error while fetching user's name(Pending Invites)> ", err)
			(*project.PendingInvites)[i].Name = (*project.PendingInvites)[i].Email[:strings.IndexByte((*project.PendingInvites)[i].Email, '@')]
		}
	}

	projectID, err := project.CreateProject()
	if err != nil {
		log.Println("CreateProjectHandler: ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("failed to create project"))
		return

	}

	responses.JSON(w, http.StatusCreated, projectID)

}

// GetProjectsHandler is the handler for /getprojects route
func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	projects, err := project.GetProjects(tokenStruct.Claims["email"].(string))
	if err != nil {
		log.Println("GetProjectsHandler: ", err)
		responses.ERROR(w, http.StatusBadRequest, errors.New("failed to fetch projects"))
	}
	responses.JSON(w, http.StatusOK, projects)
}

// AddProjectMembersHandler creates handler for /addnewprojectmembers route
func AddProjectMembersHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	userDetails, _ := utils.GetUserDetails(tokenStruct.UID)

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Println("AddProjectMembersHandler: <Error while decoding request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	for i := 0; i < len(*project.PendingInvites); i++ {
		(*project.PendingInvites)[i].InvitedBy = userDetails.DisplayName
		(*project.PendingInvites)[i].Name, err = utils.GetUserDetailsByEmail((*project.PendingInvites)[i].Email)
		if err != nil {
			log.Println("AddProjectMembersHandler: <Error while fetching user's name(pending Invites) ", err)
			(*project.PendingInvites)[i].Name = (*project.PendingInvites)[i].Email[:strings.IndexByte((*project.PendingInvites)[i].Email, '@')]
		}
	}

	err = project.AddProjectMembers(tokenStruct.Claims["email"].(string))
	if err != nil {
		log.Println("AddProjectMembersHandler: ", err)
		responses.ERROR(w, http.StatusUnauthorized, errors.New("failed to add new member"))
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
	userDetails, err := utils.GetUserDetails(tokenStruct.UID)
	if err != nil {
		log.Println("AcceptInviteHandler: <Error while fetching user details> ", err)
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Println("AcceptInviteHandler: <Error while decoding request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = project.AcceptInvite(userDetails)
	if err != nil {
		log.Println("AcceptInviteHandler: ", err)
		responses.ERROR(w, http.StatusBadRequest, errors.New("failed to accept invitation"))
		return
	}
	responses.JSON(w, http.StatusOK, "Project Invitation Accepted")
}

func RejectInviteHandler(w http.ResponseWriter, r *http.Request) {
	project := &models.Project{}
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	userDetals, _ := utils.GetUserDetails(tokenStruct.UID) // todo: handle error

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Println("RejectInviteHandler: <failed to decode request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = project.RejectInvite(userDetals)
	if err != nil {
		log.Println("RejectInviteHandler: ", err)
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
		log.Println("LeaveProjectHandler: <failed to decode request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = project.LeaveProject(tokenStruct.Claims["email"].(string))
	if err != nil {
		log.Println("LeaveProjectHandler: ", err)
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
		ProjectID string `json:"projectID,omitempty"`
		Email     string `json:"email,omitempty"`
	}
	var detail *Details
	err := json.NewDecoder(r.Body).Decode(&detail)
	if err != nil {
		log.Println("RemoveProjectMemberHandler: <failed to decode request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = models.RemoveProjectMember(tokenStruct.Claims["email"].(string), detail.Email, detail.ProjectID)
	if err != nil {
		log.Println("RemoveProjectMemberHandler: ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, nil)
}

func GetPendingInvitesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)

	invites, err := models.GetPendingInvites(tokenStruct.Claims["email"].(string))

	if err != nil {
		log.Println("GetPendingInvitesHandler: ", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	responses.JSON(w, http.StatusOK, invites)
}

func GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	members, err := models.TeamMembers(vars["projectID"])
	if err != nil {
		log.Println("GetTeamMembersHandler: ", err)
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	responses.JSON(w, http.StatusOK, members)
}

func ChangeMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := ctx.Value("tokenuid")
	tokenStruct := token.(*auth.Token)
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Println("ChangeMemberRoleHandler: <failed to decode request body> ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	err = models.ChangeMemberRole(requestBody, tokenStruct.Claims["email"].(string))
	if err != nil {
		log.Println("ChangeMemberRoleHandler: ", err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	} else {
		responses.JSON(w, http.StatusOK, "Updated user role")
	}
}
