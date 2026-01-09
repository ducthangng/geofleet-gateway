package apis

import (
	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
)

type RideHandler struct {
	ride_v1.UnimplementedRideServiceServer
}

func NewRideHandler() *RideHandler {
	return &RideHandler{}
}

func (rhl *RideHandler) TrackMultipleRides(input *ride_v1.TrackMultipleRidesRequest, stream ride_v1.RideService_TrackMultipleRidesServer) (err error) {
	return
}
