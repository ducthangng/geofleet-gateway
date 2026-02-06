package apis

import (
	"context"
	"errors"
	"io"
	"log"

	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/google/uuid"
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

// Risk: Bottleneck when there are a lot of ride request
func (rhl *RideHandler) RequestRide(ctx context.Context, input *ride_v1.RequestRideRequest) (*ride_v1.RequestRideResponse, error) {
	result, err := rhl.RideClient.RequestRide(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rhl *RideHandler) SubscribeRideUpdate(input *ride_v1.SubscribeRideUpdateRequest, stream grpc.ServerStreamingServer[ride_v1.SubscribeRideUpdateResponse]) error {

	// a stream from the user --> gateway --> rideservice
	ctx := stream.Context()

	rideServiceStream, err := rhl.RideClient.SubscribeRideUpdate(ctx, input)
	if err != nil {
		return nil
	}

	for {
		msg, err := rideServiceStream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Println("rideServiceStream_e1: ", err)
			return err
		}

		// when receive the message, send back immediately to the user
		if err = stream.SendMsg(msg); err != nil {
			log.Println("rideServiceStream_e2: ", err)
			return err
		}
	}
}

func (rhl *RideHandler) TrackMultipleRides(input *ride_v1.TrackMultipleRidesRequest, stream ride_v1.RideService_TrackMultipleRidesServer) (err error) {
	return
}

// Driver accept the ride
func (rhl *RideHandler) AcceptRide(ctx context.Context, input *ride_v1.AcceptRideRequest) (*ride_v1.AcceptRideResponse, error) {
	isValidUUID := rhl.isValidUUID(ctx, input.Value.GetRideId())
	if !isValidUUID {
		return nil, errors.New("invalid ride id - CODE: SR526")
	}

	// forward the request
	resp, err := rhl.RideClient.AcceptRide(ctx, input)
	if err != nil {
		return nil, errors.New("accept ride failed. - CODE: SR526")
	}

	return resp, nil
}

func (rhl *RideHandler) StartRide(ctx context.Context, input *ride_v1.StartRideRequest) (*ride_v1.StartRideResponse, error) {
	isValidUUID := rhl.isValidUUID(ctx, input.Value.GetRideId())
	if !isValidUUID {
		return nil, errors.New("invalid ride id - CODE: SR526")
	}

	resp, err := rhl.RideClient.StartRide(ctx, input)
	if err != nil {
		return nil, errors.New("start ride failed. - CODE: SR526")
	}

	return resp, nil
}

// User / Driver cancel the ride
func (rhl *RideHandler) CancelRide(ctx context.Context, input *ride_v1.CancelRideRequest) (*ride_v1.CancelRideResponse, error) {
	isValidUUID := rhl.isValidUUID(ctx, input.Value.GetRideId())
	if !isValidUUID {
		return nil, errors.New("invalid ride id - CODE: SR526")
	}

	if len(input.Reason) == 0 {
		return nil, errors.New("invalid reason - CODE: INVALID")
	}

	// do not allow to cancel if the user are currently in the ride
	// must complete
	resp, err := rhl.RideClient.CancelRide(ctx, input)
	if err != nil {
		return nil, errors.New("start ride failed. - CODE: SR526")
	}

	return resp, nil
}

func (rhl *RideHandler) CompleteRide(ctx context.Context, input *ride_v1.CompleteRideRequest) (*ride_v1.CompleteRideResponse, error) {

	isValidUUID := rhl.isValidUUID(ctx, input.Value.GetRideId())
	if !isValidUUID {
		return nil, errors.New("invalid ride id - CODE: SR526")
	}

	resp, err := rhl.RideClient.CompleteRide(ctx, input)
	if err != nil {
		return nil, errors.New("start ride failed. - CODE: SR526")
	}

	return resp, nil
}

func (rhl *RideHandler) isValidUUID(ctx context.Context, id string) bool {
	if len(id) == 0 {
		return false
	}

	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	return true
}
