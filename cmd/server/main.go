package main

import (
	"log"
	"net/http"

	"metrics-tpl/internal/server"
	"metrics-tpl/internal/server/storage"
)

func main() {
	config, err := server.GetConfig()
	if err != nil {
		log.Println(err)
	}

	storage := storage.New()

	router, err := server.NewRouter(storage)
	if err != nil {
		log.Println(err)
	}

	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Println(err)
	}
}
