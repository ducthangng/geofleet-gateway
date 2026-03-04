package kakfa_consumer

import (
	"context"
	"io"
	"log"

	"github.com/bytedance/sonic"
	"github.com/ducthangng/geofleet/gateway/app/domain/entity"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaManager struct {
	kafkaReader *kafka.Reader
}

func NewKafkaManager() *KafkaManager {
	return &KafkaManager{
		kafkaReader: singleton.GetKafkaReader(singleton.KAFKA_TOPIC_UPDATE_RIDES),
	}
}

type RideUpdatesParams struct {
	RideID uuid.UUID
}

func (k *KafkaManager) ListenRideUpdates(ctx context.Context) (<-chan entity.Ride, error) {

	// allow read simutaneously
	readOnlyC := make(chan entity.Ride, 100)

	go func() {

		// when this go-routine is broken, channel close immediately, no leak
		defer close(readOnlyC)
		for {
			var (
				err     error
				ride    entity.Ride
				message kafka.Message
			)

			message, err = k.kafkaReader.ReadMessage(ctx)
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Print("error reading kafka message: ", err)
			}

			if err = sonic.Unmarshal(message.Value, &ride); err != nil {
				log.Print("error reading kafka message: ", err)
				continue
			}

			readOnlyC <- ride
		}
	}()

	return readOnlyC, nil
}
