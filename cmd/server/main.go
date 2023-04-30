package main

import (
	"log"
	"net/http"

	"metrics-tpl/internal/server"
	"metrics-tpl/internal/server/storage"
)

func main() {
	storage := storage.New()

	router, err := server.NewRouter(storage)
	if err != nil {
		log.Println(err)
	}

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println(err)
	}
}
