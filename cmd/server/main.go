package main

import (
	"log"
	"metrics-tpl/internal/server"
	"metrics-tpl/internal/server/storage"
	"net/http"
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
