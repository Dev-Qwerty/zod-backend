package models

import (
	"context"
	"log"
	"time"

	"firebase.google.com/go/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
)

// User model
type User struct {
	ID string `pg:"id,notnull,unique" json:"id"`
	// Fname     string    `pg:"notnull" json:"fname"`
	// Lname     string    `pg:"notnull" json:"lname"`
	Name      string    `pg:"name,notnull" json:"name"`
	Email     string    `pg:"email,notnull,unique" json:"email"`
	CreatedAt time.Time `pg:"createdat,notnull" json:"created_at"`
	Project   []string  `pg:"projects,array" json:"projects"`
}

// CreateNewUser creates a new user
func CreateNewUser() {
	params := (&auth.UserToCreate{}).
		DisplayName("").
		Email("").
		Password("").
		EmailVerified(false)

	u, err := config.Client.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		log.Println("user created")
		user := &User{}

		user.ID = u.UID
		user.Name = u.DisplayName
		user.Email = u.Email
		user.CreatedAt = time.Now()

		_, err := database.DB.Model(user).Insert()
		if err != nil {
			log.Printf("Failed to insert user to db: %v", err)
		} else {
			log.Println("user inserted")
		}
	}
}
