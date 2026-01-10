package apis

import (
	"io"
	"log"
	"time"

	"github.com/mmcloughlin/geohash"

	"github.com/bytedance/sonic"
	tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/segmentio/kafka-go"
)

type TrackingHandler struct {
	tracking_v1.UnimplementedTrackingServiceServer
}

func NewTrackingHandler() *TrackingHandler {
	return &TrackingHandler{}
}

type DriverLocationEvent struct {
	DriverID  string  `json:"driver_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp int64   `json:"timestamp"`
	Geohash   string  `json:"geohash"`
}

func (thl *TrackingHandler) UploadLocationHistory(stream tracking_v1.TrackingService_UploadLocationHistoryServer) error {
	ctx := stream.Context()
	var count int32

	log.Println("UploadLocationHistory....")
	log.Println("starting kafka")
	kafkaWriter := singleton.GetKafkaWriter()
	log.Println("done starting kafka")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("end stream: %v", err)

			// Khi kết thúc stream, trả về tổng số điểm đã xử lý
			return stream.SendAndClose(&tracking_v1.UploadLocationHistoryResponse{
				Msg:             "Data streamed to Kafka successfully",
				PointsProcessed: count,
			})
		}

		log.Println("read successfully")
		if err != nil {
			log.Printf("Failed to read: %v", err)
			return err
		}

		// 1. Tính toán Geohash từ tọa độ Lat/Lng
		// Độ chính xác (precision) 5 hoặc 6 là chuẩn cho quy mô thành phố (khoảng 1km - 4km)
		hash := geohash.Encode(req.Location.Lat, req.Location.Lng)
		geoKey := hash[:5] // Lấy 5 ký tự đầu để nhóm các xe trong cùng 1 khu vực vào 1 Partition

		event := DriverLocationEvent{
			DriverID:  req.DriverId,
			Lat:       req.Location.Lat,
			Lng:       req.Location.Lng,
			Timestamp: time.Now().Unix(),
			Geohash:   hash, // hash đầy đủ
		}

		log.Println("marshalling...")
		payload, err := sonic.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshall: %v", err)
			return err
		}

		log.Println("writing message ...")
		err = kafkaWriter.WriteMessages(ctx, kafka.Message{
			Topic: "telemetry.driver.locations",
			Key:   []byte(geoKey), //  DriverID is the key
			Value: payload,
		})

		if err != nil {
			log.Printf("Failed to push to Kafka: %v", err)
			// Optional: can continue depend on the type of data
			// here, the data is a stream of locaiton, which mean it does not matter if one didnt receive
			continue
		}

		log.Println("increased count")
		count++
	}
}
