package models

import (
	"context"
	"log"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
)

// User model
type User struct {
	ID        string    `pg:"id,notnull,unique" json:"id"`
	Fname     string    `pg:"fname,notnull" json:"fname"`
	Lname     string    `pg:"lname,notnull" json:"lname"`
	Email     string    `pg:"email,notnull,unique" json:"email"`
	CreatedAt time.Time `pg:"createdat,notnull" json:"created_at"`
	Project   []string  `pg:"projects,array" json:"projects"`
}

// FirebaseUser model
type FirebaseUser struct {
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateNewUser creates a new user
func CreateNewUser(user FirebaseUser) {
	displayName := user.Fname + " " + user.Lname
	// Save user to firebase
	params := (&auth.UserToCreate{}).
		DisplayName(displayName).
		Email(user.Email).
		Password(user.Password).
		EmailVerified(false)

	u, err := config.Client.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		log.Println("user created")
		pguser := &User{}

		pguser.ID = u.UID
		pguser.Fname = user.Fname
		pguser.Lname = user.Lname
		pguser.Email = user.Email
		pguser.CreatedAt = time.Now()

		// Save user to postgres
		_, err := database.DB.Model(pguser).Insert()
		if err != nil {
			log.Printf("Failed to insert user to db: %v", err)
		} else {
			log.Println("user inserted")
		}
	}
}
