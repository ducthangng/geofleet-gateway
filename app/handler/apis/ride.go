package apis

import (
	"io"

	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"google.golang.org/grpc"
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
func (rhl *RideHandler) RequestRide(input *ride_v1.RequestRideRequest, stream grpc.ServerStreamingServer[ride_v1.RequestRideResponse]) error {
	ctx := stream.Context()

	rideStream, err := rhl.RideClient.RequestRide(ctx, &ride_v1.RequestRideRequest{})
	if err != nil {
		return err
	}

	for {
		// call ride service
		resp, err := rideStream.Recv()
		if err == io.EOF {
			// end of file, return nil
			return nil
		}

		if err != nil {
			return err
		}

		stream.Send(resp)
	}
}

func (rhl *RideHandler) AcceptRide() {

}

func (rhl *RideHandler) CancelRide() {

}
