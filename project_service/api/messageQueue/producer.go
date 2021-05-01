package messageQueue

import (
	"context"
	"encoding/json"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func WriteMessage(key string, value interface{}) {
	message, err := json.Marshal(value)
	if err != nil {
		log.Println("KafkaWriteMessage: ", err)
	}
	err = KafakWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: message,
		},
	)

	if err != nil {
		log.Println("KafkaProducer <write error>: ", err)
	}
}
