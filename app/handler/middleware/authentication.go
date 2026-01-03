package middleware

import (
	"log"
	"net/http"

	"github.com/ducthangng/geofleet/gateway/service/gjwt"
	"github.com/gin-gonic/gin"
)

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
