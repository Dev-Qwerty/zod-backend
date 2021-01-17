package api

import (
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/gorilla/mux"
)

// Run starts the server
func Run() {

	var err error

	router := mux.NewRouter()

	// Adds firebase to the server
	err = config.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize firebase: %v", err)
	}

	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start with error: %v", err)
	}
}
