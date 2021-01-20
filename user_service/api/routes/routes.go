package routes

import (
	"github.com/Dev-Qwerty/zod-backend/user_service/api/controllers"
	"github.com/gorilla/mux"
)

// InitializeRoutes init routes
func InitializeRoutes(r *mux.Router) {
	s := r.PathPrefix("/api/user").Subrouter()
	s.HandleFunc("/signup", controllers.UserSignUp).Methods("POST")
}
