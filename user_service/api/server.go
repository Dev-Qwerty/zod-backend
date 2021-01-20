package api

import (
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/routes"
	"github.com/gorilla/mux"
)

// Run starts the server
func Run() {

	var err error

	router := mux.NewRouter()
	routes.InitializeRoutes(router)

	// Adds firebase to the server
	err = config.InitializeFirebase()
	if err != nil {
		log.Printf("Failed to initialize firebase: %v", err)
	}

	// Connect to postgres
	err = database.Connect()
	if err != nil {
		log.Printf("Failed to connect to db: %v", err)
	}

	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start with error: %v", err)
	}
}
