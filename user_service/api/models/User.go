package models

import (
	"context"
	"errors"
	"log"
	"math/rand"
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
	Imgurl    string    `pg:"imgurl" json:"imgurl"`
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

// Create profile avatars(UI Avatar)
func generateProfileAvatar(fname, lname string) string {
	rand.Seed(time.Now().Unix())
	colorsArray := []string{"F44336", "9C27B0", "CDDC39", "FF9800", "757575", "00ACC1", "E91E63", "004D40", "FFEB3B", "607D8B0", "4E342E", "E64A19", "4CAF50", "039BE5", "F57F17", "424242", "00BFA5", "D81B60", "006064", "5E35B1"}
	bgColor := colorsArray[rand.Intn(len(colorsArray))]
	imgUrl := "https://ui-avatars.com/api/?name=" + fname + "+" + lname + "&background=" + bgColor + "&color=fff"
	return imgUrl
}

// CreateNewUser creates a new user
func CreateNewUser(user FirebaseUser) error {
	displayName := user.Fname + " " + user.Lname

	imgUrl := generateProfileAvatar(user.Fname, user.Fname)

	// Save user to firebase
	params := (&auth.UserToCreate{}).
		DisplayName(displayName).
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
	newuser.Fname = user.Fname
	newuser.Lname = user.Lname
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

	err = utils.SendEmailVerificationLink(user.Fname, user.Email)
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
