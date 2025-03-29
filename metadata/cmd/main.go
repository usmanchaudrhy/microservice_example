package main

import (
	"awesomeProject/metadata/internal/controller/metadata"
	httphandler "awesomeProject/metadata/internal/handler/http"
	"awesomeProject/metadata/internal/repository/memory"
	"awesomeProject/pkg/discovery"
	"awesomeProject/pkg/discovery/consul"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const serviceName = "metadata"

func main() {

	var port int
	flag.IntVar(&port, "port", 8081, "Api handler port")
	flag.Parse()

	log.Printf("starting metadata service on port %d", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
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

	log.Println("starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)

	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
