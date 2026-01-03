package singleton

import (
	"sync"

	pb "github.com/ducthangng/geofleet-proto/user"
	_ "github.com/mbobakov/grpc-consul-resolver" // IMPORTANT: REGISTER consul resolver with grpc
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient pb.UserServiceClient
	userConn   *grpc.ClientConn
	userOnce   sync.Once
)

// GetUserServiceClient trả về singleton client để gọi sang User Service
func GetUserServiceClient() (pb.UserServiceClient, error) {
	var err error

	userOnce.Do(func() {

		// URI format: consul://[consul-addr]/[service-name]?wait=14s
		// wait=14s giúp thực hiện "long polling", subscribe thay đổi gần như realtime
		target := "consul://127.0.0.1:8500/user-service?wait=14s"

		// Khởi tạo kết nối duy nhất một lần
		userConn, err = grpc.NewClient(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// Bạn có thể thêm các cấu trúc như Keepalive để duy trì kết nối
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		)
		if err == nil {
			userClient = pb.NewUserServiceClient(userConn)
		}
	})

	return userClient, err
}

// Hàm này dùng để đóng kết nối khi Gateway tắt (Graceful Shutdown)
func CloseUserConn() {
	if userConn != nil {
		userConn.Close()
	}
}
