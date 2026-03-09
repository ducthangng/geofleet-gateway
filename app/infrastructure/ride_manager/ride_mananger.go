package ride_manager

import (
	"context"
	"log"
	"sync"

	"github.com/ducthangng/geofleet/gateway/app/domain/entity"
	"github.com/ducthangng/geofleet/gateway/app/infrastructure/kakfa_consumer"
	"github.com/ducthangng/geofleet/gateway/app/infrastructure/redis_consumer"
	"github.com/google/uuid"
)

/*
RideManager manage the information of ride updates via Kafka.
It operates as a doubled-linked-list that can be evicted if the ride's is either:
  - End (finish, cancelled)
  - Are not update recently (in the nearest 10.000 rides)
*/
type RideManager struct {
	kafka *kakfa_consumer.KafkaManager
	redis *redis_consumer.RedisConsumer
	mutex sync.RWMutex
	store map[string]*Ride
}

var rideManager *RideManager

type Ride struct {
	Data      entity.Ride
	ID        string
	Location  string                      // Current Lat/Lng
	Observers map[string]chan entity.Ride // Key: SessionID, Value: The Channel
	mu        sync.RWMutex                // Protects the Observers map
}

func Initialize() *RideManager {
	rideManager = &RideManager{
		kafka: kakfa_consumer.NewKafkaManager(),
		redis: redis_consumer.NewRedisConsumer(),
		store: make(map[string]*Ride),
	}

	return rideManager
}

func GetRideManager() *RideManager {
	return rideManager
}

func (r *RideManager) UpdateRideStatus(ctx context.Context) error {
	rideUpateChan, err := r.kafka.ListenRideUpdates(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case update, ok := <-rideUpateChan:
			if !ok {
				// channel is closed
				return nil
			}

			// got the update of ride, start update into redis and into cache
			// update redis
			if err = r.redis.UpdateRide(ctx, update); err != nil {
				log.Println("update rides return error: ", err)
			}

			// update cache, implement later

			// notify the new-update to observer of the ride
			if err = r.Notify(update); err != nil {
				log.Println("notify return error: ", err)
				return err
			}
		}
	}
}

func (r *RideManager) Notify(data entity.Ride) error {
	// when accessing the store, you must safely lock the data so that mutext does not affect the
	r.mutex.Lock()
	defer r.mutex.Unlock()

	storedRide := r.store[data.RideID.String()]
	// locked stored ride, so that another one if accessing need to wait. This might happen if another subcriber is accessing it ...
	storedRide.mu.Lock()
	defer storedRide.mu.Unlock()

	for i := range storedRide.Observers {
		// notify
		ch := storedRide.Observers[i]
		ch <- data
	}

	return nil
}

// When the updates from Kafka arrived, get the ride, update it
// then, return the channel data to the caller
func (r *RideManager) Subscribe(ctx context.Context, rideID string) (string, chan entity.Ride) {
	r.mutex.Lock()
	storedRide, ok := r.store[rideID]
	if !ok {
		// initialize
		storedRide = &Ride{
			ID:        rideID,
			Observers: map[string]chan entity.Ride{},
		}

		r.store[rideID] = storedRide
	}
	defer r.mutex.Unlock()

	storedRide.mu.Lock()
	defer storedRide.mu.Unlock()

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	// IMPORTANT: Use a buffered channel (size 1) to avoid blocking Kafka
	ch := make(chan entity.Ride, 1)
	storedRide.Observers[sessionID.String()] = ch

	// TODO: returns all the previous data to the channel

	return sessionID.String(), ch
}

func (r *RideManager) Unsubscribe(rideID string, sessionID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	storedRide := r.store[rideID]

	storedRide.mu.Lock()
	defer storedRide.mu.Unlock()

	delete(storedRide.Observers, sessionID)

	return nil
}
