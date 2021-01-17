package config

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Client gets an auth client from the firebase.App
var Client *auth.Client

// InitializeFirebase adds firebase to the server
func InitializeFirebase() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}

	firebaseCred := os.Getenv("FIREBASE_CREDENTIALS")

	sa := option.WithCredentialsFile(firebaseCred)
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return err
	}

	Client, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}
