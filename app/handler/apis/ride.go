package apis

import (
	"context"
	"errors"
	"log"
	"time"

	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	"github.com/ducthangng/geofleet/gateway/app/infrastructure/ride_manager"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RideHandler struct {
	ride_v1.UnimplementedRideServiceServer
	RideClient  ride_v1.RideServiceClient
	RideManager *ride_manager.RideManager
}

func NewRideHandler() *RideHandler {
	return &RideHandler{
		RideClient:  singleton.NewRideClient(),
		RideManager: ride_manager.Initialize(),
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

// this can be called directly inside request ride
func (rhl *RideHandler) TriggerRideUpdate(ctx context.Context, input *ride_v1.TriggerRideUpdateRequest) (*ride_v1.TriggerRideUpdateResponse, error) {
	response, err := rhl.RideClient.TriggerRideUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SubscribeRideUpdate: the client will receive the stream of response from the service via this method.
func (rhl *RideHandler) SubscribeRideUpdate(input *ride_v1.SubscribeRideUpdateRequest, stream grpc.ServerStreamingServer[ride_v1.SubscribeRideUpdateResponse]) error {
	ctx := stream.Context()

	// subscribe to updates
	sessionID, ch := rhl.RideManager.Subscribe(ctx, input.RideId)
	if len(sessionID) != 36 {
		log.Println("error occur at subscribing: ", sessionID)
	}

	previousUpdate := []*ride_v1.RideStatusUpdateResponse{}

	for {
		select {
		case <-ctx.Done():
			// Unsubscribe the data
			rhl.RideManager.Unsubscribe(input.RideId, sessionID)

		case data := <-ch:
			if err := stream.Send(&ride_v1.SubscribeRideUpdateResponse{
				RideId: input.RideId,
				RideStatus: &ride_v1.RideStatusUpdateResponse{
					RideStatus:  ride_v1.RideStatus(data.RideStatus),
					LastUpdated: timestamppb.New(time.Now()),
				},
				PreviousRideStatus: previousUpdate,
				LastUpdated:        timestamppb.New(time.Now()),
			}); err != nil {
				return err
			}

			previousUpdate = append(previousUpdate, &ride_v1.RideStatusUpdateResponse{
				RideStatus:  ride_v1.RideStatus(data.RideStatus),
				LastUpdated: timestamppb.New(time.Now()),
			})
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
