package main

import (
	"awesomeProject/rating/internal/controller/rating"
	httphandler "awesomeProject/rating/internal/handler/http"
	"awesomeProject/rating/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
