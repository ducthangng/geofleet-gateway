package singleton

import (
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaOnce   sync.Once
)

// GetKafkaWriter returns a singleton instance of the Kafka Producer
func GetKafkaWriter(brokers []string, topic string) *kafka.Writer {
	kafkaOnce.Do(func() {
		kafkaWriter = &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{}, // Efficiently distributes messages
			BatchTimeout: 10 * time.Millisecond,
			Async:        true, // Non-blocking for Gateway performance
		}
	})
	return kafkaWriter
}

// CloseKafka cleans up the connection on shutdown
func CloseKafka() error {
	if kafkaWriter != nil {
		return kafkaWriter.Close()
	}
	return nil
}
