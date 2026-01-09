package apis

import (
	tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"
)

type TrackingHandler struct {
	tracking_v1.UnimplementedTrackingServiceServer
}

func NewTrackingHandler() *TrackingHandler {
	return &TrackingHandler{}
}

func (thl *TrackingHandler) UploadLocationHistory(input tracking_v1.TrackingService_UploadLocationHistoryServer) error {
	return nil
}
