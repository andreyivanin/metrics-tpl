package main

import (
	"log"
	"metrics-tpl/internal/server"
	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/storage"
	"net/http"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New(config)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.ApplyConfig()
	if err != nil {
		log.Fatal(err)
	}

	router, err := server.NewRouter(storage)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
