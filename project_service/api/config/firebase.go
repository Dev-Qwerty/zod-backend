package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client is the interface for the Firebase auth service
var Client *auth.Client

// InitializeFirebase creates a new Firebase auth client
func InitializeFirebase() error {
	type param struct {
		Type         string `json:"type"`
		Project_ID   string `json:"project_id"`
		Private_Key  string `json:"private_key"`
		Client_Email string `json:"client_email"`
	}

	firebase_param := param{
		Type:         os.Getenv("TYPE"),
		Project_ID:   os.Getenv("PROJECT_ID"),
		Private_Key:  os.Getenv("PRIVATE_KEY"),
		Client_Email: os.Getenv("CLIENT_EMAIL"),
	}

	firebase_opt, err := json.Marshal(firebase_param)

	if err != nil {
		log.Println("marshall error: ", err)
	}

	opt := option.WithCredentialsJSON([]byte(firebase_opt))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	Client, err = app.Auth(context.TODO())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return nil
}
