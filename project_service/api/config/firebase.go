package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client is the interface for the Firebase auth service
var Client *auth.Client

// InitializeFirebase creates a new Firebase auth client
func InitializeFirebase() error {
	opt := option.WithCredentialsFile("/home/albin/ServiceAccountKey-zode.json")
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