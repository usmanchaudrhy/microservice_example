package main

import (
	"awesomeProject/metadata/internal/controller/metadata"
	httphandler "awesomeProject/metadata/internal/handler/http"
	"awesomeProject/metadata/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)

	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
