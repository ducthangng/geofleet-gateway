package repository

import (
	"context"

	"github.com/ducthangng/geofleet/gateway/app/domain/entity"
)

type KafkaRepository interface {
	ListenRideUpdates(ctx context.Context) (<-chan entity.Ride, error)
}
