package config

import (
	"flag"
	"log"
	"time"

	// "time"

	"github.com/caarlos0/env/v6"
)

const (
	SERVERADDRPORT = "localhost:8080"
	STOREINTERVAL  = 300
	STOREFILE      = "metrics-db.json"
	RESTORE        = true
)

type Config struct {
	Address          string        `env:"ADDRESS"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL"`
	StoreFile        string        `env:"STORE_FILE"`
	RestoreSavedData bool          `env:"RESTORE"`
}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address and port")
	StoreIntervalFlag := flag.Int("i", STOREINTERVAL, "server store interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "server db store file")
	flag.BoolVar(&cfg.RestoreSavedData, "r", cfg.RestoreSavedData, "server restore db from file on start?")

	flag.Parse()

	StoreIntervalDuration := time.Duration(*StoreIntervalFlag) * time.Second
	cfg.StoreInterval = StoreIntervalDuration

	return nil
}

func getEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func Read() (Config, error) {
	var cfg = Config{
		Address:          SERVERADDRPORT,
		StoreInterval:    STOREINTERVAL * time.Second,
		StoreFile:        STOREFILE,
		RestoreSavedData: RESTORE,
	}

	getFlag(&cfg)
	getEnv(&cfg)

	return cfg, nil

}
