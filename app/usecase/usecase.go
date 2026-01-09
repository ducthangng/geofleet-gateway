package usecase

import (
	"context"

	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service"
	"google.golang.org/grpc"
)

type IdentityServiceUsecase interface {
	CheckDuplicatedPhone(ctx context.Context, dto user_service.CheckDuplicatePhoneRequest) (user_service.CheckDuplicatePhoneResponse, error)
	CreateUserProfile(ctx context.Context, data user_service.UserCreation) (userId string, err error)
	Login(ctx context.Context, in user_service.LoginRequest, opts ...grpc.CallOption) (user_service.LoginResponse, error)
	GetUserProfile(ctx context.Context, userId string) (res user_service.User, err error)
}
