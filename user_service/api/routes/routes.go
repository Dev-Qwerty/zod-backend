package routes

import (
	"github.com/Dev-Qwerty/zod-backend/user_service/api/controllers"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/middlewares"
	"github.com/gorilla/mux"
)

// InitializeRoutes init routes
func InitializeRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/api/user").Subrouter()            //User routes
	projectRouter := r.PathPrefix("/api/user/project").Subrouter() //Project routes
	etm := r.PathPrefix("/api/user").Subrouter()                   //Extract token middleware

	userRouter.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	projectRouter.HandleFunc("/new", controllers.NewProject).Methods("POST")
	projectRouter.HandleFunc("/update", controllers.UpdateProject).Methods("PUT")
	projectRouter.HandleFunc("/delete", controllers.DeleteProject).Methods("DELETE")
	etm.Use(middlewares.ExtractToken)
	etm.HandleFunc("/update", controllers.Update).Methods("PUT")
	etm.HandleFunc("/delete", controllers.Delete).Methods("DELETE")
}
