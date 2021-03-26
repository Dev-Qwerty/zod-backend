package models

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/Dev-Qwerty/zod-backend/project_service/api/database"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Project model
type Project struct {
	ProjectID      string           `json:"projectID,omitempty" bson:"_id,omitempty"`
	ProjectName    string           `json:"projectName,omitempty" bson:"projectName,omitempty"`
	Members        *[]Member        `json:"projectMembers,omitempty" bson:"projectMembers,omitempty"`
	PendingInvites *[]PendingInvite `json:"pendingInvites,omitempty" bson:"pendingInvites,omitempty"`
	Teamlead       string           `json:"teamlead,omitempty" bson:"teamlead,omitempty"`
	Deadline       string           `json:"deadline,omitempty" bson:"deadline,omitempty"`
}

// Member model
type Member struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	UserID   string `json:"userID,omitempty" bson:"userID,omitempty"`
	MemberID string `json:"memberID,omitempty" bson:"memberID,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Role     string `json:"userRole,omitempty" bson:"userRole,omitempty"`
}

// PendingInvite models
type PendingInvite struct {
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	Role      string `json:"userRole,omitempty" bson:"userRole,omitempty"`
	InvitedBy string `json:"invitedby,omitempty" bson:"invitedby,omitempty"`
}

// CreateProject creates a new project and save it to db
func (p *Project) CreateProject() (string, error) {
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	p.ProjectID = ksuid.New().String()
	createdProjectID, err := zodeProjectCollection.InsertOne(context.TODO(), p)
	if err != nil {
		return "", err
	}
	id := createdProjectID.InsertedID.(string)
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
		return []*Project{}, err
	}
	for cursor.Next(context.TODO()) {
		var project Project
		err := cursor.Decode(&project)
		if err != nil {
			return []*Project{}, err
		}
		projects = append(projects, &project)
	}

	cursor.Close(context.TODO())

	if err := cursor.Err(); err != nil {
		return []*Project{}, err
	}

	return projects, nil

}

func (p *Project) AddProjectMembers(email string) error {
	userEmail := email
	filter := bson.M{
		"_id": p.ProjectID,
		"projectMembers": bson.M{
			"$elemMatch": bson.M{
				"email":    userEmail,
				"userRole": "Owner",
			},
		},
	}
	var project *Project
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	err := zodeProjectCollection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		return err
	}
	filter = bson.M{"_id": p.ProjectID}
	update := bson.M{"$push": bson.M{"pendingInvites": bson.M{"$each": p.PendingInvites}}}
	_, err = zodeProjectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) AcceptInvite(userDetails *auth.UserInfo) error {
	var project *Project
	userEmail := userDetails.Email
	filter := bson.M{
		"_id":                  p.ProjectID,
		"pendingInvites.email": userEmail,
	}
	projection := bson.M{
		"pendingInvites.$": 1,
	}
	// fetch role of user from db
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	err := zodeProjectCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&project)
	if err != nil {
		return err
	}
	pendinginvite := &[]PendingInvite{}
	pendinginvite = project.PendingInvites

	member := &Member{}
	member.Name = userDetails.DisplayName
	member.Email = userEmail
	member.UserID = userDetails.UID
	member.Role = (*pendinginvite)[0].Role

	// add user to project members
	filter = bson.M{"_id": p.ProjectID}
	update := bson.M{"$push": bson.M{"projectMembers": member}}
	_, err = zodeProjectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	update = bson.M{
		"$pull": bson.M{
			"pendingInvites": bson.M{
				"email": userEmail,
			},
		}}
	_, err = zodeProjectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) LeaveProject(email string) error {
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	// stages for Aggregate pipeline
	matchStage := bson.D{
		primitive.E{Key: "$match", Value: bson.M{
			"_id": p.ProjectID,
		}},
	}
	memberCount := bson.D{
		primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$size", Value: "$projectMembers"},
			}},
		}},
	}

	//count total no of members in the project
	cursor, err := zodeProjectCollection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, memberCount})
	if err != nil {
		return err
	}

	var result bson.M
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return err
		}
	}
	cursor.Close(context.TODO())
	// delete project if there is only one member
	if result["count"] == int32(1) {
		// todo: delete chat and boards in the project
		_, err = zodeProjectCollection.DeleteOne(context.TODO(), bson.M{"_id": p.ProjectID})
		if err != nil {
			return err
		}
		return nil
	}
	// remove member if there is more than one member in the project
	filter := bson.M{
		"_id":                  p.ProjectID,
		"projectMembers.email": email,
	}
	update := bson.M{
		"$pull": bson.M{
			"projectMembers": bson.M{
				"email": email,
			},
		},
	}
	_, err = zodeProjectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func RemoveProjectMember(email, memberEmail string, projectID primitive.ObjectID) error {
	userEmail := email
	// check if the user is owner
	filter := bson.M{
		"_id": projectID,
		"projectMembers": bson.M{
			"$elemMatch": bson.M{
				"email":    userEmail,
				"userRole": "Owner",
			},
		},
	}
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	update := bson.M{
		"$pull": bson.M{
			"projectMembers": bson.M{
				"email": memberEmail,
			},
		}}
	result, err := zodeProjectCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return err
	}
	return nil
}

func GetPendingInvites(email string) ([]map[string]interface{}, error) {
	fmt.Println("get invites")
	var invites []map[string]interface{}
	projection := bson.M{
		"projectName":               1,
		"teamlead":                  1,
		"pendingInvites.userRole.$": 1,
	}

	filter := bson.M{
		"pendingInvites": bson.M{
			"$elemMatch": bson.M{
				"email": email,
			},
		},
	}
	fmt.Println(email)
	zodeProjectCollection := database.Client.Database("zodeProjectDB").Collection("projects")
	cursor, err := zodeProjectCollection.Find(context.TODO(), filter, options.Find().SetProjection(projection))

	if err != nil {
		fmt.Println(err)
		return invites, err
	}
	var invite map[string]interface{}
	for cursor.Next(context.TODO()) {

		err := cursor.Decode(&invite)
		if err != nil {
			fmt.Println(err)
			return invites, err
		}
		invites = append(invites, invite)
	}

	cursor.Close(context.TODO())

	if err := cursor.Err(); err != nil {
		return invites, err
	}

	return invites, nil
}
