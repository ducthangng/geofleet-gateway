package singleton

import (
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaOnce   sync.Once
)

// GetKafkaWriter returns a singleton instance of the Kafka Producer
func GetKafkaWriter() *kafka.Writer {
	kafkaOnce.Do(func() {
		centralConfig := GetGlobalConfig()
		brokers := centralConfig.KafkaBrokers

		kafkaWriter = &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{}, // Efficiently distributes messages
			BatchTimeout: 10 * time.Millisecond,
			Async:        true, // Non-blocking for Gateway performance
			Logger:       kafka.LoggerFunc(log.Printf),
			ErrorLogger:  kafka.LoggerFunc(log.Printf),
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
