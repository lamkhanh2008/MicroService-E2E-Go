package main

import (
	"log"
	metadatagateway "microservice/movie/gateway/metadata/http"
	ratinggateway "microservice/movie/gateway/rating/http"
	"microservice/movie/internal/controller/movie"
	httphandler "microservice/movie/internal/handler/http"
	"net/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
