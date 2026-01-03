package gjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/ducthangng/geofleet/gateway/app/singleton"
	"github.com/golang-jwt/jwt/v5"
)

type JWTEncodingType struct {
	UserId    string
	Role      string
	SessionId string
	Phone     string
}

type SigningClaims struct {
	Data JWTEncodingType
	jwt.RegisteredClaims
}

func EncodeJWT(data JWTEncodingType) (string, error) {
	jwtSecretKey := singleton.GlobalConfig.JwtSecretKey

	// encode jwt
	newClaims := SigningClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-auth-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*SigningClaims, error) {
	jwtSecretKey := singleton.GlobalConfig.JwtSecretKey

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &SigningClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Check that the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*SigningClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
