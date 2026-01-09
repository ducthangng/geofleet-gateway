package middleware

import (
	"log"
	"net/http"

	"context"
	"strings"

	"github.com/ducthangng/geofleet/gateway/service/gjwt"
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

		if publicMethods[info.FullMethod] {
			// bypassing unary interceptor
			return handler(ctx, req)
		}

		// Checking authorization credentials
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
	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return userId, status.Errorf(codes.Unauthenticated, "invalid auth header format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 4. Validate JWT
	claims, err := gjwt.VerifyToken(tokenString)
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

		decodedSigningKey, err := gjwt.VerifyToken(cookie)
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
		ctx.Set("SessionID", decodedSigningKey.Data.SessionId)
	}
}
