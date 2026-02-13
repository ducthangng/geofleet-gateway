package ride_manager

import (
	"sync"

	"github.com/ducthangng/geofleet/gateway/app/domain/entity"
)

type RideManager struct {
	mutex sync.RWMutex
	store map[string]*entity.Ride
}
