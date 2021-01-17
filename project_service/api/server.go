package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Run creates and starts server
func Run() {
	router := mux.NewRouter()

	log.Println("server starting on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start with err: %v", err)
	}
}
