package messageQueues

import "github.com/segmentio/kafka-go"

const (
	topic         = "testproject"
	brokerAddress = "localhost:9092"
)

var (
	kw *kafka.Writer
	kr *kafka.Reader
)

func InitializeKafka() {
	kw = &kafka.Writer{
		Addr: kafka.TCP(brokerAddress),
		Topic: topic,
	}

	kr = kafka.NewReader(kafka.ReaderConfig{
		Brokers:                []string{brokerAddress},
		Topic:                  topic,
		Partition: 0,
		GroupID: "group1",
		StartOffset:            kafka.LastOffset,
	})
}
