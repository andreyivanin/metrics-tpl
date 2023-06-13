package main

import (
	"context"
	"log"
	"metrics-tpl/internal/server"
	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/storage"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	router := server.NewRouter(storage)

	log.Printf("Running http server on port: %s", config.Address)

	err = http.ListenAndServe(config.Address, router)
	if err != nil {
		log.Fatal(err)
	}
}
