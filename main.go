package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ducthangng/geofleet/gateway/app/handler"
	"github.com/ducthangng/geofleet/gateway/app/singleton"
)

func main() {
	singleton.InitializeConfig()
	centralConfig := singleton.GetGlobalConfig()

	singleton.GetRedisClient()
	singleton.GetConsulClient()
	singleton.GetKafkaWriter(centralConfig.KafkaBrokers, centralConfig.KafkaTopic)

	handlr := handler.Routing()

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", centralConfig.Host, centralConfig.Port),
		ReadTimeout:  time.Duration(centralConfig.RequestTimeout * int(time.Second)),
		WriteTimeout: time.Duration(centralConfig.RequestTimeout * int(time.Second)),
		IdleTimeout:  time.Duration(centralConfig.RequestTimeout * int(time.Second)),
		Handler:      handlr,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
