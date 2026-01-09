package middleware

import (
	"context"

	"google.golang.org/grpc"
)

// WrappedStream giúp chúng ta ghi đè Context của stream cũ bằng Context mới chứa userID
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func AuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// (1) Kiểm tra Whitelist cho stream
		if publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		// (2) Authorize
		userID, err := authorize(ss.Context())
		if err != nil {
			return err
		}

		// (3) Gắn userID vào context và bọc stream lại
		newCtx := context.WithValue(ss.Context(), UserIDKey, userID)
		return handler(srv, &wrappedStream{ServerStream: ss, ctx: newCtx})
	}
}
