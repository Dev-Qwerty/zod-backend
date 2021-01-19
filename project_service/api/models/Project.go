package models

// Project model
type Project struct {
	// ProjectID   string     `json:"projectID,omitempty" bson:"projectID,omitempty"`
	ProjectName    string           `json:"projectName,omitempty" bson:"projectName,omitempty"`
	Channels       *[]Channel       `json:"channels,omitempty" bson:"channels,omitempty"`
	Members        *[]Member        `json:"projectMembers,omitempty" bson:"projectMembers,omitempty"`
	PendingInvites *[]PendingInvite `json:"pendingInvites,omitempty" bson:"pendingInvites,omitempty"`
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
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Role   string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}
