package config

import (
	"flag"
	"log"
	"os"
	"strconv"
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
	Address          string `env:"ADDRESS"`
	StoreInterval    time.Duration
	StoreFile        string `env:"STORE_FILE"`
	RestoreSavedData bool   `env:"RESTORE"`
}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address and port")
	StoreIntervalFlag := flag.Int("i", STOREINTERVAL, "server store interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "server db store file")
	flag.BoolVar(&cfg.RestoreSavedData, "r", cfg.RestoreSavedData, "server restore db from file on start?")

	flag.Parse()

	cfg.StoreInterval = time.Duration(*StoreIntervalFlag) * time.Second

	return nil
}

func getEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		log.Print(err)
		return err
	}

	StoreIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if !ok {
		return nil
	}

	StoreIntervalEnvInt, err := strconv.Atoi(StoreIntervalEnv)
	if err != nil {
		return err
	}

	cfg.StoreInterval = time.Duration(StoreIntervalEnvInt) * time.Second

	return nil
}

func Read() (Config, error) {
	var err error

	var cfg = Config{
		Address:          SERVERADDRPORT,
		StoreInterval:    STOREINTERVAL * time.Second,
		StoreFile:        STOREFILE,
		RestoreSavedData: RESTORE,
	}

	err = getFlag(&cfg)
	if err != nil {
		return Config{}, err
	}

	err = getEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil

}
