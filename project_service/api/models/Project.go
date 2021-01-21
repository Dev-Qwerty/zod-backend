package models

import (
	"context"
	"fmt"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Project model
type Project struct {
	ProjectID      string           `json:"projectID,omitempty" bson:"_id,omitempty"`
	ProjectName    string           `json:"projectName,omitempty" bson:"projectName,omitempty"`
	Channels       *[]Channel       `json:"channels,omitempty" bson:"channels,omitempty"`
	Members        *[]Member        `json:"projectMembers,omitempty" bson:"projectMembers,omitempty"`
	PendingInvites *[]PendingInvite `json:"pendingInvites,omitempty" bson:"pendingInvites,omitempty"`
	Teamlead       string           `json:"teamlead,omitempty" bson:"teamlead,omitempty"`
	Deadline       string           `json:"deadline,omitempty" bson:"deadline,omitempty"`
}

// Member model
type Member struct {
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Role   string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}

// PendingInvite models
type PendingInvite struct {
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Role  string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}

// CreateProject creates a new project and save it to db
func (p *Project) CreateProject() (string, error) {
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	createdProjectID, err := zodeProjectCollection.InsertOne(context.TODO(), p)
	if err != nil {
		return "", err
	}
	id := createdProjectID.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// GetProjects fetch projects of a user from mongodb
func (p *Project) GetProjects(userEmail string) ([]*Project, error) {
	var projects []*Project
	projection := bson.M{
		"projectName": 1,
		"teamlead":    1,
		"deadline":    1,
	}
	filter := bson.M{
		"projectMembers": bson.M{
			"$elemMatch": bson.M{
				"email": userEmail}}}
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	cursor, err := zodeProjectCollection.Find(context.TODO(), filter, options.Find().SetProjection(projection))

	if err != nil {
		fmt.Println("Error fetching projects: ", err)
	}
	for cursor.Next(context.TODO()) {
		var project Project
		err := cursor.Decode(&project)
		if err != nil {
			fmt.Println("error decoding docs: ", err)
		}
		projects = append(projects, &project)
	}

	cursor.Close(context.TODO())

	if err := cursor.Err(); err != nil {
		return []*Project{}, err
	}

	return projects, nil

}
