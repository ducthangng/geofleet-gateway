package user_service

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ducthangng/geofleet-proto/user"
	"github.com/ducthangng/geofleet/gateway/app/singleton"

	"google.golang.org/grpc"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

type UserCreation struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Bod      string `json:"bod"`
}

func (u *UserService) CreateUserProfile(ctx context.Context, data UserCreation) (userId int, err error) {

	log.Println("enter CreateUserProfile")
	conn, err := singleton.GetUserServiceClient()
	if err != nil {
		log.Println("error 1: ", err)
		return
	}

	grpcRes, err := conn.CreateUserProfile(ctx, &pb.UserCreationRequest{
		Fullname: data.Fullname,
		Email:    data.Email,
		Phone:    data.Phone,
		Password: data.Password,
		Address:  data.Address,
		Bod:      data.Bod,
	})
	if err != nil {
		log.Println("error 2: ", err)
		return userId, err
	}

	log.Println("got this: ", grpcRes)

	return
}

type UserDTO struct {
	UserId   int     `json:"userId"`
	Fullname string  `json:"fullname"`
	Address  string  `json:"address"`
	Bod      string  `json:"bod"`
	Score    float64 `json:"score"`
}

// this is a unary call, call once and get the result immediately
func (u *UserService) GetUserProfile(ctx context.Context, userId int) (res UserDTO, err error) {
	conn, err := singleton.GetUserServiceClient()
	if err != nil {
		return res, err
	}

	grpcResult, err := conn.GetUserProfile(ctx, &pb.UserInfoRequest{
		UserId: fmt.Sprintf("%d", userId),
	})
	if err != nil {
		return res, err
	}

	res = UserDTO{
		UserId:   userId,
		Fullname: grpcResult.Fullname,
		Address:  grpcResult.Address,
		Bod:      grpcResult.Bod,
		Score:    grpcResult.Score,
	}

	return res, nil
}

func (u *UserService) Login(ctx context.Context, data UserDTO) (user UserDTO, err error) {
	return
}

func (u *UserService) TrackMultipleRides(*pb.TrackRidesRequest, grpc.ServerStreamingServer[pb.RideLocation]) error {

	return nil
}

func (u *UserService) UploadLocationHistory(grpc.ClientStreamingServer[pb.LocationData, pb.UploadStatus]) error {
	return nil
}
