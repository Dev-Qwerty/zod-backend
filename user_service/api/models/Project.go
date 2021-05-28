package models

import (
	"log"
	"strings"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/database"
)

// NewProject model
type Project struct {
	ProjectId string        `json:"projectid"`
	Member    ProjectMember `json:"member"`
}

type ProjectMember struct {
	Name   string `json:"name"`
	Fid    string `json:"fid"`
	ImgUrl string `json:"imgUrl"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}

// Add project of user to db
func AddProject(project Project) error {
	_, err := database.DB.Exec(`UPDATE users SET projects = array_append(projects, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)
	if strings.ToLower(project.Member.Role) == "admin" {
		_, err = database.DB.Exec(`UPDATE users SET owner = array_append(owner, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)
	}
	if err != nil {
		log.Printf("Error at AddProject User.go : %v", err)
		return err
	}

	return nil
}

// Update project role of users
func UpdateProject(project Project) error {
	if strings.ToLower(project.Member.Role) == "admin" {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_append(owner, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)

		if err != nil {
			log.Printf("Error at UpdateProject User.go : %v", err)
			return err
		}
	} else {
		_, err := database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)

		if err != nil {
			log.Printf("Error at UpdateProject User.go : %v", err)
			return err
		}
	}

	return nil
}

// Delete project of users
func DeleteProject(project Project) error {

	result, err := database.DB.Exec(`SELECT * FROM users WHERE id = '` + project.Member.Fid + `' AND '` + project.ProjectId + `' = ANY(owner);`)
	if err != nil {
		log.Printf("Error at DeleteProject User.go : %v", err)
		return err
	}

	if result.RowsReturned() == 1 {
		_, err = database.DB.Exec(`UPDATE users SET owner = array_remove(owner, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)

		if err != nil {
			log.Printf("Error at DeleteProject User.go : %v", err)
			return err
		}
	}

	_, err = database.DB.Exec(`UPDATE users SET projects = array_remove(projects, '` + project.ProjectId + `') WHERE id = '` + project.Member.Fid + `';`)

	if err != nil {
		log.Printf("Error at DeleteProject User.go : %v", err)
		return err
	}

	return nil
}
