package user_service

import (
	"context"

	common_v1 "github.com/ducthangng/geofleet-proto/gen/go/common/v1"
	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"

	"github.com/ducthangng/geofleet/gateway/app/singleton"

	"google.golang.org/grpc"
)

type IdentityService struct {
	identity_v1.UnimplementedUserServiceServer
}

func NewUserService() *IdentityService {
	return &IdentityService{}
}

/*
Task: passing the request to the other service
*/
func (u *IdentityService) CreateUserProfile(ctx context.Context, data UserCreation) (userId string, err error) {
	conn, err := singleton.GetUserServiceClient()
	if err != nil {
		return
	}

	grpcRes, err := conn.CreateUserProfile(ctx, &identity_v1.CreateUserProfileRequest{
		Fullname: data.Fullname,
		Email:    data.Email,
		Phone:    data.Phone,
		Password: &common_v1.Password{
			Value: data.Password,
		},
		Address: data.Address,
		// TODO: convert string - timestamp to timestamppb
		// Bod:      data.Bod,
	})

	if err != nil {
		return userId, err
	}

	return grpcRes.UserId, nil
}

// this is a unary call, call once and get the result immediately
func (u *IdentityService) GetUserProfile(ctx context.Context, userId string) (res User, err error) {
	conn, err := singleton.GetUserServiceClient()
	if err != nil {
		return res, err
	}

	grpcResult, err := conn.GetUserProfile(ctx, &identity_v1.GetUserProfileRequest{
		UserId: userId,
	})
	if err != nil {
		return res, err
	}

	res = User{
		UserId:   userId,
		Fullname: grpcResult.Fullname,
		Address:  grpcResult.Address,
		// TODO: what is the good practice here?
		// Bod:   *grpcResult.Bod,
		Score: grpcResult.Score,
	}

	return res, nil
}

func (u *IdentityService) Login(context.Context, LoginRequest, ...grpc.CallOption) (LoginResponse, error) {
	return LoginResponse{}, nil
}

func (u *IdentityService) CheckDuplicatedPhone(ctx context.Context, data CheckDuplicatePhoneRequest) (CheckDuplicatePhoneResponse, error) {
	return CheckDuplicatePhoneResponse{}, nil
}
