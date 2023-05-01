package server

import (
	"flag"
	"log"

	// "time"

	"github.com/caarlos0/env/v6"
)

const (
	SERVERADDRPORT = "localhost:8080"
	// STOREINTERVAL  = 10
	// STOREFILE      = "devops-metrics-db.json"
	// RESTORE        = true
)

type Config struct {
	Address string `env:"ADDRESS"`
	// StoreInterval    int `env:"STORE_INTERVAL"`
	// StoreFile        string        `env:"STORE_FILE"`
	// RestoreSavedData bool          `env:"RESTORE"`
}

func GetFlagConfig(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address and port")
	// flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "server store interval")
	// flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "server db store file")
	// flag.BoolVar(&cfg.RestoreSavedData, "r", cfg.RestoreSavedData, "server restore db from file on start?")
	flag.Parse()
	return nil
}

func GetEnvConfig(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetConfig() (Config, error) {
	var cfg = Config{
		Address: SERVERADDRPORT,
		// StoreInterval:    STOREINTERVAL,
		// StoreFile:        STOREFILE,
		// RestoreSavedData: RESTORE,
	}

	GetFlagConfig(&cfg)
	GetEnvConfig(&cfg)

	return cfg, nil

}
