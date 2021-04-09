package routes

import (
	"github.com/Dev-Qwerty/zod-backend/user_service/api/controllers"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/middlewares"
	"github.com/gorilla/mux"
)

// InitializeRoutes init routes
func InitializeRoutes(r *mux.Router) {

	r.Use(middlewares.CorsHandler)

	userRouter := r.PathPrefix("/api/user").Subrouter()            //User routes
	projectRouter := r.PathPrefix("/api/user/project").Subrouter() //Project routes
	etm := r.PathPrefix("/api/user").Subrouter()                   //Extract token middleware

	userRouter.HandleFunc("/signup", controllers.SignUp).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/resendverificationemail", controllers.ResendEmail).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/resetpassword", controllers.ResetPassword).Methods("POST", "OPTIONS")
	projectRouter.HandleFunc("/new", controllers.NewProject).Methods("POST", "OPTIONS")
	projectRouter.HandleFunc("/update", controllers.UpdateProject).Methods("PUT", "OPTIONS")
	projectRouter.HandleFunc("/delete", controllers.DeleteProject).Methods("DELETE", "OPTIONS")
	etm.Use(middlewares.ExtractToken)
	etm.HandleFunc("/update", controllers.Update).Methods("PUT", "OPTIONS")
	etm.HandleFunc("/delete", controllers.Delete).Methods("DELETE", "OPTIONS")
}
