package singleton

import (
	"log"
	"sync"

	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rideClientOnce sync.Once
	rideClient     ride_v1.RideServiceClient
	rideConn       *grpc.ClientConn
)

func NewRideClient() ride_v1.RideServiceClient {
	var err error

	rideClientOnce.Do(func() {

		target := "consul://127.0.0.1:8500/ride-service?wait=14s"

		rideConn, err = grpc.NewClient(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// Bạn có thể thêm các cấu trúc như Keepalive để duy trì kết nối
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		)

		if err != nil {
			log.Println("failed to connect to consul client: Ride Client")
			return
		}

		rideClient = ride_v1.NewRideServiceClient(rideConn)
	})

	return rideClient
}
