package singleton

import (
	"sync"

	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"
	tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"
	_ "github.com/mbobakov/grpc-consul-resolver" // IMPORTANT: REGISTER consul resolver with grpc
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	trackingClient tracking_v1.TrackingServiceClient
	trackingConn   *grpc.ClientConn
	trackingOnce   sync.Once
)

// GetUserServiceClient trả về singleton client để gọi sang User Service
func GetTrackingClient() (tracking_v1.TrackingServiceClient, error) {
	var err error

	userOnce.Do(func() {

		// URI format: consul://[consul-addr]/[service-name]?wait=14s
		// wait=14s giúp thực hiện "long polling", subscribe thay đổi gần như realtime
		target := "consul://127.0.0.1:8500/tracking-service?wait=14s"

		// Khởi tạo kết nối duy nhất một lần
		userConn, err = grpc.NewClient(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// Bạn có thể thêm các cấu trúc như Keepalive để duy trì kết nối
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		)
		if err == nil {
			userClient = identity_v1.NewUserServiceClient(userConn)
		}
	})

	return trackingClient, err
}

// Hàm này dùng để đóng kết nối khi Gateway tắt (Graceful Shutdown)
func CloseTrackingConn() {
	if userConn != nil {
		userConn.Close()
	}
}
