package api

import (
	"log"
	"net/http"
	"os"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/database"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Run creates and starts server
func Run() {
	router := mux.NewRouter()

	routes.InitializeRoutes(router)

	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load environment variables: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
		log.Println("server starting on port 8080")

		port := ":" + os.Getenv("PORT")

		err := http.ListenAndServe(port, router)
		if err != nil {
			log.Fatalf("Server failed to start with err: %v", err)
		}
	}()

	if err := config.InitializeFirebase(); err != nil {
		log.Printf("Failed to initialize firebase: %v", err)
	}
	database.InitializeDB()

}
