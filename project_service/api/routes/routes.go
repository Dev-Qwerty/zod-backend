package routes

import (
	"github.com/Dev-Qwerty/zod-backend/project_service/api/controllers"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/middlewares"
	"github.com/gorilla/mux"
)

// InitializeRoutes initializes routes
func InitializeRoutes(r *mux.Router) {
	r.Use(middlewares.Cors)
	r.Use(middlewares.ExtractUID)
	s := r.PathPrefix("/api/projects").Subrouter()
	s.HandleFunc("/createproject", controllers.CreateProjectHandler).Methods("POST", "OPTIONS")
	s.HandleFunc("/getprojects", controllers.GetProjectsHandler).Methods("GET", "OPTIONS")
	s.HandleFunc("/addnewprojectmembers", controllers.AddProjectMembersHandler).Methods("POST", "OPTIONS")
	s.HandleFunc("/acceptInvite", controllers.AcceptInviteHandler).Methods("PUT", "OPTIONS")
	s.HandleFunc("/rejectinvite", controllers.RejectInviteHandler).Methods("PUT", "OPTIONS")
	s.HandleFunc("/leaveproject", controllers.LeaveProjectHandler).Methods("PUT", "OPTIONS")
	s.HandleFunc("/removemember", controllers.RemoveProjectMemberHandler).Methods("PUT", "OPTIONS")
	s.HandleFunc("/invites", controllers.GetPendingInvitesHandler).Methods("GET", "OPTIONS")
	s.HandleFunc("/{projectID}/members", controllers.GetTeamMembers).Methods("GET", "OPTIONS")
	s.HandleFunc("/member/changerole", controllers.ChangeMemberRoleHandler).Methods("PUT", "OPTIONS")
}
