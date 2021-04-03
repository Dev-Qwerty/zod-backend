package config

import (
	"context"
	"encoding/json"
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
		Type                        string `json:"type"`
		Project_ID                  string `json:"project_id"`
		Private_Key_Id              string `json:"private_key_id"`
		Private_Key                 string `json:"private_key"`
		Client_Email                string `json:"client_email"`
		Client_Id                   string `json:"client_id"`
		Auth_Uri                    string `json:"auth_uri"`
		Token_Uri                   string `json:"token_uri"`
		Auth_Provider_x509_Cert_Url string `json:"auth_provider_x509_cert_url"`
		Client_x509_Cert_Url        string `json:"client_x509_cert_url"`
	}
	// privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")
	// r := strings.NewReplacer("\\\n", "\\n")
	// newPrivateKey := r.Replace(os.Getenv(privateKey))

	firebaseParam := param{
		Type:                        os.Getenv("FIREBASE_TYPE"),
		Project_ID:                  os.Getenv("FIREBASE_PROJECT_ID"),
		Private_Key_Id:              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		Private_Key:                 os.Getenv("FIREBASE_PRIVATE_KEY_TEST"),
		Client_Email:                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		Client_Id:                   os.Getenv("FIREBASE_CLIENT_ID"),
		Auth_Uri:                    os.Getenv("FIREBASE_AUTH_URI"),
		Token_Uri:                   os.Getenv("FIREBASE_TOKEN_URI"),
		Auth_Provider_x509_Cert_Url: os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		Client_x509_Cert_Url:        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}

	firebaseCred, err := json.Marshal(firebaseParam)
	if err != nil {
		return err
	}

	sa := option.WithCredentialsJSON([]byte(firebaseCred))
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
