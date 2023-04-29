package main

import (
	"log"
	"net/http"

	"metrics-tpl/internal/server/handler"
	"metrics-tpl/internal/server/storage"
)

func main() {
	storage := storage.New()
	customHandler := handler.NewHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/", customHandler.MetricUpdate)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err)
	}
}
