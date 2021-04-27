package models

import (
	"context"
	"errors"
	"log"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/utils"
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
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Email string `json:"email"`
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
		log.Printf("Error at CreateNewUser User.go : %v", err)
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
		log.Printf("Error at CreateNewUser User.go : %v", err)
		return err
	}

	err = utils.SendEmailVerificationLink(user.Fname ,user.Email)
	if err != nil {
		log.Printf("Error at CreateNewUser User.go : %v", err)
		return err
	}

	return nil
}

// FetchUser fetch a single user details from db
func FetchUser(uid string) (map[string]string, error) {
	user := &User{}
	resp := make(map[string] string)
	err := database.DB.Model(user).Column("fname", "lname", "email").Where("id = ?", uid).Select()
	if err != nil {
		log.Printf("Error at FetchUser User.go: %v", err)
		return resp, err
	}

	resp["fname"] = user.Fname
	resp["lname"] = user.Lname
	resp["email"] = user.Email
	
	return resp, nil
}

// UpdateUser updates the user data
func UpdateUser(data UpdatedUser) error {
	user := &User{}
	user.ID = data.ID

	if data.Email != "" {
		user.Email = data.Email
		_, err := database.DB.Model(user).Set("email = ?email").Where("id = ?id").Update()
		if err != nil {
			log.Printf("Error at UpdatUser User.go : %v", err)
			return err
		}
	}
	if data.Fname != "" {
		user.Fname = data.Fname
		_, err := database.DB.Model(user).Set("fname = ?fname").Where("id = ?id").Update()
		if err != nil {
			log.Printf("Error at UpdatUser User.go : %v", err)
			return err
		}
	}
	if data.Lname != "" {
		user.Lname = data.Lname
		_, err := database.DB.Model(user).Set("lname = ?lname").Where("id = ?id").Update()
		if err != nil {
			log.Printf("Error at UpdatUser User.go : %v", err)
			return err
		}
	}

	return nil
}

// DeleteUser removes the user from firebase and postgres
func DeleteUser(id string) error {
	user := &User{}
	user.ID = id
	err := database.DB.Model(user).Where("id = ?id").Select()
	if err != nil {
		log.Printf("Error at DeleteUser User.go : %v", err)
		return err
	}
	if len(user.Role) == 0 {
		err := config.Client.DeleteUser(context.Background(), id)

		if err != nil {
			log.Printf("Error at DeleteUser User.go : %v", err)
			return err
		}
		_, err = database.DB.Model(user).Where("id = ?id").Delete()

		if err != nil {
			log.Printf("Error at DeleteUser User.go : %v", err)
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
		log.Printf("Error at AddProject User.go : %v", err)
		return err
	}

	return nil
}

// Update project role of users
func UpdateProject(project Project) error {
	if project.Role == "owner" {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_append(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			log.Printf("Error at UpdateProject User.go : %v", err)
			return err
		}
	} else {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			log.Printf("Error at UpdateProject User.go : %v", err)
			return err
		}
	}

	return nil
}

// Delete project of users
func DeleteProject(project Project) error {

	result, err := database.DB.Exec(`SELECT * FROM users WHERE id = '` + project.ID + `' AND '` + project.ProjectId + `' = ANY(owner);`)
	if err != nil {
		log.Printf("Error at DeleteProject User.go : %v", err)
		return err
	}

	if result.RowsReturned() == 1 {
		_, err = database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

		if err != nil {
			log.Printf("Error at DeleteProject User.go : %v", err)
			return err
		}
	}

	_, err = database.DB.Exec(`UPDATE users SET projects = array_remove(projects, '` + project.ProjectId + `') WHERE id = '` + project.ID + `';`)

	if err != nil {
		log.Printf("Error at DeleteProject User.go : %v", err)
		return err
	}

	return nil
}
