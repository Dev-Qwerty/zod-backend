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
	Role      []string  `pg:"owner,array" json:"role"`
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
	ID    string
	Field string
	Value string
}

// NewProject model
type Project struct {
	ID        string
	ProjectId string
	Role      string
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

	newuser := &User{}

	newuser.ID = u.UID
	newuser.Fname = user.Fname
	newuser.Lname = user.Lname
	newuser.Email = user.Email
	newuser.CreatedAt = time.Now()

	// Save user to postgres
	_, err = database.DB.Model(newuser).Insert()
	if err != nil {
		config.Client.DeleteUser(context.Background(), u.UID)
		return err
	}

	return nil
}

// UpdateUser updates the user data
func UpdateUser(data UpdatedUser) error {
	user := &User{}
	user.ID = data.ID

	switch data.Field {
	case "email":
		user.Email = data.Value
		_, err := database.DB.Model(user).Set("email = ?email").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	case "fname":
		user.Fname = data.Value
		_, err := database.DB.Model(user).Set("fname = ?fname").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	case "lname":
		user.Lname = data.Value
		_, err := database.DB.Model(user).Set("lname = ?lname").Where("id = ?id").Update()
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid field")
	}
	return nil
}

// DeleteUser removes the user from firebase and postgres
func DeleteUser(id string) error {
	user := &User{}
	user.ID = id
	err := database.DB.Model(user).Where("id = ?id").Select()
	if err != nil {
		return err
	}
	if len(user.Role) == 0 {
		err := config.Client.DeleteUser(context.Background(), id)

		if err != nil {
			return err
		}
		_, err = database.DB.Model(user).Where("id = ?id").Delete()

		if err != nil {
			return err
		}
	} else {
		return errors.New("User is owner of some projetcs")
	}

	return nil
}

// Add project of user to db
func AddProject(project Project) error {
	_, err := database.DB.Exec(`UPDATE users SET projects = array_append(projects, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)
	if project.Role == "owner" {
		_, err = database.DB.Exec(`UPDATE users SET owner = array_append(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)
	}
	if err != nil {
		return err
	}

	return nil
}

// Update project role of users
func UpdateProject(project Project) error {
	if project.Role == "owner" {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_append(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			return err
		}
	} else {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			return err
		}
	}

	return nil
}

// Delete project of users
func DeleteProject(project Project) error {

	result, err := database.DB.Exec(`SELECT * FROM users WHERE id = '` + project.ID + `' AND '` + project.ProjectId + `' = ANY(owner);`)
	if err != nil {
		return err
	}

	if result.RowsReturned() == 1 {
		_, err = database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			return err
		}
	}

	_, err = database.DB.Exec(`UPDATE users SET projects = array_remove(projects, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

	if err != nil {
		return err
	}

	return nil
}
