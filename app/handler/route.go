package handler

import (
	"github.com/ducthangng/geofleet/gateway/app/handler/restful"
	"github.com/gin-gonic/gin"
)

func Routing() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	globalHandler := restful.NewGlobalHandler()
	r.POST("/api/register", globalHandler.Register)

	// userHandler := restful.NewUserHandler()

	// v1 := r.Group("/user")
	// v1.GET("/", userHandler.GetUser) // get user info
	// v1.GET("/past_rides")
	// v1.POST("/login")
	// v1.POST("/online")

	// v2 := r.Group("/ride")
	// v2.GET("/current_activity")
	// v2.POST("/prepare")
	// v2.POST("/match")
	// v2.POST("/end")
	// v2.POST("/rate")

	return r
}
