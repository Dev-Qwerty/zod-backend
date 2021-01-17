package api

import (
	"log"
	"net/http"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/database"
	"github.com/gorilla/mux"
)

// Run creates and starts server
func Run() {
	router := mux.NewRouter()

	defer func() {
		if r := recover(); r != nil {
			log.Println("DB connection failed: ", r)
		}
		log.Println("server starting on port 8080")
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatalf("Server failed to start with err: %v", err)
		}
	}()

	database.InitializeDB()

}
