package usecase

import (
	"context"

	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service"
	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service/pb"
	"google.golang.org/grpc"
)

type UserServiceUsecase interface {
	CreateUserProfile(ctx context.Context, data user_service.UserCreation) (userId int, err error)
	GetUserProfile(ctx context.Context, userId int) (res user_service.UserDTO, err error)
	TrackMultipleRides(*pb.TrackRidesRequest, grpc.ServerStreamingServer[pb.RideLocation]) error
	UploadLocationHistory(grpc.ClientStreamingServer[pb.LocationData, pb.UploadStatus]) error
}
