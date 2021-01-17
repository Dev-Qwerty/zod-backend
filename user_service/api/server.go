package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Run starts the server
func Run() {
	router := mux.NewRouter()

	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start with error %v", err)
	}
}
