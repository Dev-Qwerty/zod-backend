package messageQueues

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/Dev-Qwerty/zod-backend/user_service/api/models"
)

func Consume() {
	msg, err := kr.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error at ReadMessage consumer.go: %v", err)
	}
	var data map[string]interface{}

	project := models.Project{}

	if len(data) != 0 {
		project.Member.Fid = data["member"].(map[string]interface{})["fid"].(string)
		project.ProjectId = data["projectId"].(string)
		project.Member.Role = data["member"].(map[string]interface{})["fid"].(string)

		err = json.Unmarshal(msg.Value, &data)
		if bytes.NewBuffer(msg.Key).String() == "Create project" {
			err = models.AddProject(project)
			if err != nil {
				log.Printf("Error at AddProject consumer.go: %v", err)
			}
		}
		if err != nil {
			log.Printf("Error at Umarshal consumer.go: %v", err)
		}
	}

}
