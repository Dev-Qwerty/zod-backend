package models

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/config"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
	"github.com/Dev-Qwerty/zod-backend/user_service/api/utils"
)

// User model
type User struct {
	ID        string    `pg:"id,notnull,unique" json:"id"`
	Name      string    `pg:"name,notnull" json:"name"`
	Email     string    `pg:"email,notnull,unique" json:"email"`
	Imgurl    string    `pg:"imgurl,notnull" json:"imgurl"`
	CreatedAt time.Time `pg:"createdat,notnull" json:"created_at"`
	Project   []string  `pg:"projects,array" json:"projects"`
	Role      []string  `pg:"owner,array" json:"role"`
}

// FirebaseUser model
type FirebaseUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdatedUser model
type UpdatedUser struct {
	ID    string
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Create profile avatars(UI Avatar)
func generateProfileAvatar(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "+")
	rand.Seed(time.Now().Unix())
	colorsArray := []string{"F44336", "9C27B0", "CDDC39", "FF9800", "757575", "00ACC1", "E91E63", "004D40", "FFEB3B", "607D8B0", "4E342E", "E64A19", "4CAF50", "039BE5", "F57F17", "424242", "00BFA5", "D81B60", "006064", "5E35B1"}
	bgColor := colorsArray[rand.Intn(len(colorsArray))]
	imgUrl := "https://ui-avatars.com/api/?name=" + name + "&background=" + bgColor + "&color=fff"
	return imgUrl
}

// CreateNewUser creates a new user
func CreateNewUser(user FirebaseUser) error {

	imgUrl := generateProfileAvatar(user.Name)

	// Save user to firebase
	params := (&auth.UserToCreate{}).
		DisplayName(user.Name).
		Email(user.Email).
		Password(user.Password).
		EmailVerified(false).
		PhotoURL(imgUrl)

	u, err := config.Client.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error at CreateNewUser User.go : %v", err)
		return err
	}

	newuser := &User{}

	newuser.ID = u.UID
	newuser.Name = user.Name
	newuser.Email = user.Email
	newuser.Imgurl = imgUrl
	newuser.CreatedAt = time.Now()

	// Save user to postgres
	_, err = database.DB.Model(newuser).Insert()
	if err != nil {
		config.Client.DeleteUser(context.Background(), u.UID)
		log.Printf("Error at CreateNewUser User.go : %v", err)
		return err
	}

	name := strings.Fields(user.Name)[0]
	err = utils.SendEmailVerificationLink(name, user.Email)
	if err != nil {
		log.Printf("Error at CreateNewUser User.go : %v", err)
		return err
	}

	return nil
}

// FetchUser fetch a single user details from db
func FetchUser(uid string) (map[string]string, error) {
	user := &User{}
	resp := make(map[string]string)
	err := database.DB.Model(user).Column("name", "email").Where("id = ?", uid).Select()
	if err != nil {
		log.Printf("Error at FetchUser User.go: %v", err)
		return resp, err
	}

	resp["name"] = user.Name
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
	if data.Name != "" {
		imgUrl := generateProfileAvatar(data.Name)
		user.Name = data.Name
		user.Imgurl = imgUrl

		params := (&auth.UserToUpdate{}).DisplayName(user.Name).PhotoURL(imgUrl)

		_, err := config.Client.UpdateUser(context.Background(), user.ID, params)
		if err != nil {
			log.Printf("Error at UpdatUser User.go : %v", err)
			return err
		}
		// _, err = database.DB.Model(user).WherePK().Update()
		_, err = database.DB.Model(user).Set("name = ?name").Where("id = ?id").Update()
		if err != nil {
			log.Printf("Error at UpdatUser User.go : %v", err)
			return err
		}
		_, err = database.DB.Model(user).Set("imgurl = ?imgurl").Where("id = ?id").Update()
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
