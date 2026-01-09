package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ducthangng/geofleet/gateway/app/handler/apis"
	"github.com/ducthangng/geofleet/gateway/app/handler/middleware"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"
	ride_v1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
	tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	singleton.InitializeConfig()
	centralConfig := singleton.GetGlobalConfig()

	singleton.GetRedisClient()
	singleton.GetConsulClient()
	singleton.GetKafkaWriter(centralConfig.KafkaBrokers, centralConfig.KafkaTopic)

	// 1. Mở port TCP
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", centralConfig.Host, centralConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Khởi tạo gRPC Server
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.AuthUnaryInterceptor(),
			// Bạn có thể thêm LoggingInterceptor hoặc RecoveryInterceptor ở đây
		),
		grpc.ChainStreamInterceptor(
			middleware.AuthStreamInterceptor(),
		),
	)

	userHandler := apis.NewIdentityHandler()
	rideHandler := apis.NewRideHandler()
	trackingHandler := apis.NewTrackingHandler()

	// register the services
	identity_v1.RegisterUserServiceServer(server, userHandler)
	ride_v1.RegisterRideServiceServer(server, rideHandler)
	tracking_v1.RegisterTrackingServiceServer(server, trackingHandler)

	// 4. Xử lý Graceful Shutdown (Hủy đăng ký khi tắt app)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Đang tắt Service...")
		// Hủy đăng ký trên Consul để Gateway không gọi vào nữa
		// (Thực hiện gọi client.Agent().ServiceDeregister(serviceID))
		server.GracefulStop()
		os.Exit(0)
	}()

	// 1. Khởi tạo Health Server
	healthServer := health.NewServer()

	// 2. Đăng ký Health Service vào gRPC Server của bạn
	healthpb.RegisterHealthServer(server, healthServer)

	// 3. Đặt trạng thái là SERVING (Đang hoạt động)
	// "user-service" ở đây phải khớp với Service Name bạn đăng ký trên Consul
	healthServer.SetServingStatus("user-service", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(server)

	// 5. Start server
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
