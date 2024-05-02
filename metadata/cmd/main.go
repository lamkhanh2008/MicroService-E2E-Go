package main

import (
	"log"
	metadata "microservice/metadata/internal/controller"
	httphandler "microservice/metadata/internal/handler/http"
	"microservice/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetaData))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
