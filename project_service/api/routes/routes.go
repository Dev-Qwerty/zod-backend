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
}
