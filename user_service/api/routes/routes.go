package routes

import (
	"github.com/Dev-Qwerty/zod-backend/user_service/api/controllers"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/middlewares"
	"github.com/gorilla/mux"
)

// InitializeRoutes init routes
func InitializeRoutes(r *mux.Router) {
	r.HandleFunc("/api/user/signup", controllers.SignUp).Methods("POST")
	etm := r.PathPrefix("/api/user/").Subrouter()
	etm.Use(middlewares.ExtractToken)
	etm.HandleFunc("/update", controllers.Update).Methods("PUT")
	etm.HandleFunc("/delete", controllers.Delete).Methods("DELETE")
}
