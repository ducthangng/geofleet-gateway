package apis

import (
	"context"

	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
)

type RideHandler struct {
	ride_v1.UnimplementedRideServiceServer
	RideClient ride_v1.RideServiceClient
}

func NewRideHandler() *RideHandler {
	return &RideHandler{
		RideClient: singleton.NewRideClient(),
	}
}

func (rhl *RideHandler) TrackMultipleRides(input *ride_v1.TrackMultipleRidesRequest, stream ride_v1.RideService_TrackMultipleRidesServer) (err error) {
	return
}

// Risk: Bottleneck when there are a lot of ride request
func (rhl *RideHandler) RequestRide(ctx context.Context, input *ride_v1.RequestRideRequest) (*ride_v1.RequestRideResponse, error) {
	result, err := rhl.RideClient.RequestRide(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rhl *RideHandler) AcceptRide() {

}

func (rhl *RideHandler) CancelRide() {

}
