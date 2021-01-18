package models

// Project model
type Project struct {
	ProjectID   string     `json:"projectID,omitempty" bson:"projectID,omitempty"`
	ProjectName string     `json:"projectName,omitempty" bson:"projectName,omitempty"`
	Channels    *[]Channel `json:"channels,omitempty" bson:"channels,omitempty"`
}

// Member model
type Member struct {
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
	Role   string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}

// PendingInvites models
type PendingInvites struct {
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
	Role   string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}
