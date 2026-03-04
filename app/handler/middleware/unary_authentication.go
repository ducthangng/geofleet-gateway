package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ducthangng/geofleet/gateway/service/gwjwt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var publicMethods = map[string]bool{
	"/geofleet.identity.v1.UserService/Login":                        true,
	"/geofleet.identity.v1.UserService/CheckDuplicatedPhone":         true,
	"/geofleet.identity.v1.UserService/CreateUserProfile":            true,
	"/geofleet.tracking.v1.TrackingService/UploadLocationHistory":    true,
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
}

// AuthInterceptor kiểm tra JWT từ metadata
func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		log.Println("checking in full methods ...")
		if publicMethods[info.FullMethod] {
			// bypassing unary interceptor
			return handler(ctx, req)
		}

		// Checking authorization credentials
		log.Println("checking authorize...")
		userId, err := authorize(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// 5. Inject user_id into context
		newCtx := context.WithValue(ctx, UserIDKey, userId)

		// 6. allow request to continue
		return handler(newCtx, req)
	}
}

func authorize(ctx context.Context) (userId string, err error) {
	// 1. Get metadata from context
	log.Println("checking metadata...")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return userId, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// 2. Header authorization
	values := md.Get("authorization")
	if len(values) == 0 {
		return userId, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// 3. "Bearer <token>"
	log.Println("bearer...")
	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return userId, status.Errorf(codes.Unauthenticated, "invalid auth header format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 4. Validate JWT
	log.Println("checking claims...")
	claims, err := gwjwt.VerifyToken(tokenString)
	if err != nil {
		return userId, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return claims.Data.UserId, nil
}

// deprecated resultful authentication middleware
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// get the cookie, decode the JWT then check if the JWT is a valid one
		cookie, err := ctx.Cookie("geofleet")
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, map[string]any{
				"message": "authentication failed",
			})

			return
		}

		decodedSigningKey, err := gwjwt.VerifyToken(cookie)
		if err != nil {
			log.Println("failed")
			return
		}

		if decodedSigningKey.Data.UserId == "" {
			log.Println("failed")
			return
		}

		// decode jwt
		ctx.Set("ID", decodedSigningKey.Data.UserId)
		ctx.Set("EntityCode", decodedSigningKey.Data.Role)
		ctx.Set("Phone", decodedSigningKey.Data.Phone)
	}
}
