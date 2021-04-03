package api

import (
	"log"
	"net/http"
	"os"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Run starts the server
func Run() {

	var err error

	router := mux.NewRouter()
	routes.InitializeRoutes(router)

	err = godotenv.Load()
	if err != nil {
		log.Printf("Failed to load env variables: %v", err)
	}

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

	port := os.Getenv("PORT")

	log.Println("Server starting on port " + port + " 🚀")
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Server failed to start with error: %v", err)
	}
}
