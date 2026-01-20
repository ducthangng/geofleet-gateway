package user_service

import (
	"context"
	"log"

	common_v1 "github.com/ducthangng/geofleet-proto/gen/go/common/v1"
	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"

	"github.com/ducthangng/geofleet/gateway/app/role"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/ducthangng/geofleet/gateway/service/gwerr"
	"github.com/ducthangng/geofleet/gateway/service/gwjwt"
)

type IdentityService struct {
	identity_v1.UnimplementedUserServiceServer
	UserClient identity_v1.UserServiceClient
}

func NewUserService() *IdentityService {
	client, err := singleton.GetUserServiceClient()
	if err != nil {
		panic(err)
	}

	return &IdentityService{
		UserClient: client,
	}
}

/*
Task: passing the request to the other service
*/
func (u *IdentityService) CreateUserProfile(ctx context.Context, data UserCreation) (userId string, err error) {

	request := &identity_v1.CreateUserProfileRequest{
		Fullname: data.Fullname,
		Email:    data.Email,
		Phone:    data.Phone,
		Password: &common_v1.Password{
			Value: data.Password,
		},
		Address: data.Address,
		Role:    data.Role,
		// TODO: convert string - timestamp to timestamppb
		// Bod:      data.Bod,
	}

	log.Println("gateway receive: ", request.Role)
	grpcRes, err := u.UserClient.CreateUserProfile(ctx, request)

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

func (u *IdentityService) Login(ctx context.Context, data LoginRequest) (LoginResponse, error) {
	var (
		res     LoginResponse
		err     error
		grpcRes *identity_v1.LoginResponse
	)

	if len(data.Password.Value) == 0 || len(data.Phone) == 0 {
		return res, gwerr.ErrInvalidInput
	}

	// call user-service
	grpcRes, err = u.UserClient.Login(ctx, &identity_v1.LoginRequest{
		Phone: data.Phone,
		Password: &common_v1.Password{
			Value: data.Password.Value,
		},
	})

	if err != nil {
		log.Println("error grpc record: ", err)
		return res, err
	}

	if len(grpcRes.User.UserId) == 0 { // authenticated
		return res, gwerr.ErrInvalidAPIKey
	}

	res.User = User{
		UserId:   grpcRes.User.UserId,
		Fullname: grpcRes.User.Fullname,
		Email:    grpcRes.User.Email,
		Phone:    grpcRes.User.Phone,
		Address:  grpcRes.User.Address,
		Role:     role.Role(grpcRes.User.Role),
	}

	res.JWTToken, err = gwjwt.EncodeJWT(gwjwt.JWTEncodingType{
		UserId: res.User.UserId,
		Phone:  res.User.Phone,
		Role:   role.Role(grpcRes.User.Role).String(),
	})

	if err != nil {
		return res, err
	}

	res.IsValid = true

	return res, nil
}

func (u *IdentityService) CheckDuplicatedPhone(ctx context.Context, data CheckDuplicatePhoneRequest) (CheckDuplicatePhoneResponse, error) {
	return CheckDuplicatePhoneResponse{}, nil
}
