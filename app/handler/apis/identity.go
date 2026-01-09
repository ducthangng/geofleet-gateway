package apis

import (
	"context"
	"log"

	"buf.build/go/protovalidate"
	"github.com/ducthangng/geofleet/gateway/app/usecase"
	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service"
	"github.com/ducthangng/geofleet/gateway/service/copier"

	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"
)

type IdentityHandler struct {
	UserUsecase usecase.IdentityServiceUsecase
	identity_v1.UnimplementedUserServiceServer
}

func NewIdentityHandler() *IdentityHandler {
	return &IdentityHandler{
		UserUsecase: user_service.NewUserService(),
	}
}

/*
Task: validate the format -> fail quick
RateLimit.
*/
func (idh *IdentityHandler) CreateUserProfile(ctx context.Context, input *identity_v1.CreateUserProfileRequest) (
	output *identity_v1.CreateUserProfileResponse, err error) {

	var (
		userCreation user_service.UserCreation
		userId       string
	)

	if err = protovalidate.Validate(input); err != nil {
		// Handle failure.
		log.Println(err)
		return
	}

	// create
	userCreation = user_service.UserCreation{
		Fullname: input.Fullname,
		Password: input.Password.Value,
		Phone:    input.Phone,
		Address:  input.Address,
		Email:    input.Email,
		Bod:      input.Bod.AsTime().GoString(),
		Role:     int(input.Role),
	}

	log.Println("got userCreation: ", userCreation)
	userId, err = idh.UserUsecase.CreateUserProfile(ctx, userCreation)
	if err != nil {
		// TODO: handle error
		log.Println("error: ", err)
		return
	}

	output = &identity_v1.CreateUserProfileResponse{
		UserId: userId,
	}

	return output, nil
}

// Login can mean get online - or log in with account info
func (idh *IdentityHandler) Login(ctx context.Context, input *identity_v1.LoginRequest) (*identity_v1.LoginResponse, error) {
	var (
		err           error
		loginRequest  user_service.LoginRequest
		loginResponse user_service.LoginResponse
		res           *identity_v1.LoginResponse
	)

	if err = protovalidate.Validate(input); err != nil {
		panic(err)
		return nil, err
	}

	copier.MustCopy(&loginRequest, input)
	loginResponse, err = idh.UserUsecase.Login(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	res = &identity_v1.LoginResponse{
		UserId: loginResponse.User.UserId,
		// Fullname: loginResponse.User.Fullname,
	}

	return res, nil
}

func (idh *IdentityHandler) GetUserProfile(context.Context, *identity_v1.GetUserProfileRequest) (*identity_v1.GetUserProfileResponse, error) {
	// interceptor

	return nil, nil
}

func (idh *IdentityHandler) CheckDuplicatedPhone(ctx context.Context, input *identity_v1.CheckDuplicatedPhoneRequest) (*identity_v1.CheckDuplicatedPhoneResponse, error) {
	var (
		err              error
		phoneCheckResult user_service.CheckDuplicatePhoneResponse
	)

	// check the duplicate phone number?
	phoneCheckResult, err = idh.UserUsecase.CheckDuplicatedPhone(ctx, user_service.CheckDuplicatePhoneRequest{
		Phone: input.Phone,
	})

	if err != nil || phoneCheckResult.IsDuplicated {
		// TODO: handle error
		log.Println("error: ", err)
		return &identity_v1.CheckDuplicatedPhoneResponse{
			IsDuplicated: true,
		}, err
	}

	return &identity_v1.CheckDuplicatedPhoneResponse{
		IsDuplicated: false,
	}, nil
}
