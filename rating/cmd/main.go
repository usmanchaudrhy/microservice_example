package main

import (
	"awesomeProject/pkg/discovery"
	"awesomeProject/pkg/discovery/consul"
	"awesomeProject/rating/internal/controller/rating"
	httphandler "awesomeProject/rating/internal/handler/http"
	"awesomeProject/rating/internal/repository/memory"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const serviceName = "rating"

func main() {

	var port int
	flag.IntVar(&port, "port", 8082, "api handler port")
	flag.Parse()

	log.Printf("starting the rating service on port %d", port)

	// The registry service is actually running on port 8500
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
				log.Println("failed to report healthy state: ", err.Error())
			}

			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	log.Println("starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
