package main

import (
	"awesomeProject/movie/internal/controller/movie"
	metadatagateway "awesomeProject/movie/internal/gateway/metadata/http"
	ratinggateway "awesomeProject/movie/internal/gateway/rating/http"
	httphandler "awesomeProject/movie/internal/handler/http"
	"awesomeProject/pkg/discovery"
	"awesomeProject/pkg/discovery/consul"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// this is the gateway service for all the other services
// connects to metadata service and the rating service

const serviceName = "movie"

func main() {

	var port int
	flag.IntVar(&port, "port", 8083, "api handler port")
	flag.Parse()

	log.Println("starting the movie service on: ", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", "localhost", port)); err != nil {
		panic(err)
	}

	go func() {

		for {

			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("failed to report helathy state", err.Error())
			}
			time.Sleep(1 * time.Second)
		}

	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)

	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
