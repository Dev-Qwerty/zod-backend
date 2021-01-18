package models

// Channel model
type Channel struct {
	ChannelID   string            `json:"channelID,omitempty" bson:"channelID,omitempty"`
	ChannelName string            `json:"channelName,omitempty" bson:"channelName,omitempty"`
	Members     *[]ChannelMembers `json:"channelMembers,omitempty" bson:"channelMembers,omitempty"`
}

// ChannelMembers model
type ChannelMembers struct {
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
}
