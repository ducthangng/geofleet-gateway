package singleton

import (
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaReader *kafka.Reader
	kafkaWOnce  sync.Once
	kafkaROnce  sync.Once
)

type KafkaTopic string

const (
	KAFKA_TOPIC_UPDATE_RIDES KafkaTopic = "ride_tracking.ride_updates"
)

func (kf KafkaTopic) String() string {
	switch kf {
	case KAFKA_TOPIC_UPDATE_RIDES:
		return "ride_tracking.ride_updates"

	default:
		return "ride_tracking.ride_updates"
	}
}

// GetKafkaWriter returns a singleton instance of the Kafka Producer
func GetKafkaWriter() *kafka.Writer {
	kafkaWOnce.Do(func() {
		centralConfig := GetGlobalConfig()
		brokers := centralConfig.KafkaBrokers

		kafkaWriter = &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{}, // Efficiently distributes messages
			BatchTimeout:           10 * time.Millisecond,
			Async:                  false, // Non-blocking for Gateway performance
			MaxAttempts:            3,
			Logger:                 kafka.LoggerFunc(log.Printf),
			ErrorLogger:            kafka.LoggerFunc(log.Printf),
			AllowAutoTopicCreation: true, //
		}
	})

	return kafkaWriter
}

func GetKafkaReader(topic KafkaTopic) *kafka.Reader {
	kafkaROnce.Do(func() {
		centralConfig := GetGlobalConfig()
		brokers := centralConfig.KafkaBrokers

		kafkaReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  "gateway-ride-update-01",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
			Topic:    topic.String(),
			// Reader should be able to auto commit
			CommitInterval: time.Second,
		})
	})

	return kafkaReader
}

// CloseKafka cleans up the connection on shutdown
func CloseKafka() error {
	if kafkaWriter != nil {
		return kafkaWriter.Close()
	}
	return nil
}
