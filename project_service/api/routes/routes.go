package routes

import (
	"github.com/Dev-Qwerty/zod-backend/project_service/api/controllers"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/middlewares"
	"github.com/gorilla/mux"
)

// InitializeRoutes initializes routes
func InitializeRoutes(r *mux.Router) {
	r.Use(middlewares.ExtractUID)
	s := r.PathPrefix("/api/projects").Subrouter()
	s.HandleFunc("/createproject", controllers.CreateProjectHandler).Methods("POST")
	s.HandleFunc("/getprojects", controllers.GetProjectsHandler).Methods("GET")
	s.HandleFunc("/addnewprojectmembers", controllers.AddProjectMembersHandler).Methods("POST")
	s.HandleFunc("/acceptInvite", controllers.AcceptInviteHandler).Methods("PUT")
	s.HandleFunc("/rejectinvite", controllers.RejectInviteHandler).Methods("PUT")
	s.HandleFunc("/leaveproject", controllers.LeaveProjectHandler).Methods("PUT")
	s.HandleFunc("/removemember", controllers.RemoveProjectMemberHandler).Methods("PUT")
	s.HandleFunc("/invites", controllers.GetPendingInvitesHandler).Methods("GET")
	s.HandleFunc("/{projectID}/teams", controllers.GetTeamMembers).Methods("GET")
	s.HandleFunc("/member/changerole", controllers.ChangeMemberRoleHandler).Methods("PUT")
}
