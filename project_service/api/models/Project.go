package models

import (
	"context"

	"github.com/Dev-Qwerty/zod-backend/project_service/api/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
