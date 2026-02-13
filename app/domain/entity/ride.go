package entity

import (
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	RideID     uuid.UUID
	UserID     uuid.UUID
	Driver     uuid.UUID
	RideStatus RideStatus
	CreatedAt  time.Time
}
