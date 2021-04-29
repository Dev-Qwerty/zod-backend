package messageQueue

import (
	kafka "github.com/segmentio/kafka-go"
)

var KafakWriter kafka.Writer

func InitializeMessageQueue() {
	KafakWriter = kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "Project",
	}
}
