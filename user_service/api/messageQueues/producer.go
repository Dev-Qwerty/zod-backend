package messageQueues

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)


func Produce(key string, value []byte) error {
	err := kw.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte(key),
			Value: value,
		},
	)

	if err != nil {
		log.Printf("Error at Kafka writeMessages producer.go: %v", err)
		return err
	}
	return nil
}
