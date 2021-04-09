package utils

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/config"
)

// GetUserDetails fetch user details from firebase
func GetUserDetails(uid string) (*auth.UserInfo, error) {
	u, err := config.Client.GetUser(context.TODO(), uid)
	if err != nil {
		log.Println("GetUserDetails: ", err)
		return nil, err
	}
	return u.UserInfo, nil
}

func GetUserDetailsByEmail(email string) (string, error) {
	u, err := config.Client.GetUserByEmail(context.TODO(), email)
	if err != nil {
		log.Println("GetUserDetailsByEmail: ", err)
		return "", err
	}
	return u.DisplayName, nil
}
