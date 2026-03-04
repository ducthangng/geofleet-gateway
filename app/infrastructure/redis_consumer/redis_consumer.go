package redis_consumer

import (
	"context"
	"log"

	"github.com/ducthangng/geofleet/gateway/app/domain/entity"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisConsumer struct {
	client *redis.Client
}

func NewRedisConsumer() *RedisConsumer {
	client, err := singleton.GetRedisClient()
	if err != nil {
		log.Fatalf("initialized redis consumer return error: ", err)
	}

	return &RedisConsumer{
		client: client,
	}
}

func (r *RedisConsumer) UpdateRide(ctx context.Context, ride entity.Ride) error {

	return nil
}

func (r *RedisConsumer) GetRide(ctx context.Context, rideID uuid.UUID) error {

	return nil
}
