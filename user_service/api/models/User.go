package models

import (
	"context"
	"errors"
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

// UpdatedUser model
type UpdatedUser struct {
	Id    string
	Field string
	Value string
}

// CreateNewUser creates a new user
func CreateNewUser(user FirebaseUser) error {
	displayName := user.Fname + " " + user.Lname
	// Save user to firebase
	params := (&auth.UserToCreate{}).
		DisplayName(displayName).
		Email(user.Email).
		Password(user.Password).
		EmailVerified(false)

	u, err := config.Client.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	pguser := &User{}

	pguser.ID = u.UID
	pguser.Fname = user.Fname
	pguser.Lname = user.Lname
	pguser.Email = user.Email
	pguser.CreatedAt = time.Now()

	// Save user to postgres
	_, err = database.DB.Model(pguser).Insert()
	if err != nil {
		config.Client.DeleteUser(context.Background(), u.UID)
		return err
	}

	return nil
}

// UpdateUser updates the user data
func UpdateUser(data UpdatedUser) error {
	pguser := &User{}
	pguser.ID = data.Id

	switch data.Field {
	case "email":
		pguser.Email = data.Value
		_, err := database.DB.Model(pguser).Set("email = ?email").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	case "fname":
		pguser.Fname = data.Value
		_, err := database.DB.Model(pguser).Set("fname = ?fname").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	case "lname":
		pguser.Lname = data.Value
		_, err := database.DB.Model(pguser).Set("lname = ?lname").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	default:
		return errors.New("Invalid field")
	}
	return nil
}
