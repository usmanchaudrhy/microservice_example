package main

import (
	"awesomeProject/movie/internal/controller/movie"
	metadatagateway "awesomeProject/movie/internal/gateway/metadata/http"
	ratinggateway "awesomeProject/movie/internal/gateway/rating/http"
	httphandler "awesomeProject/movie/internal/handler/http"
	"net/http"
)

// this is the gateway service for all the other services
// connects to metadata service and the rating service

func main() {
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)

	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
