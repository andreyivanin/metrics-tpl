package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v8"
)

const (
	_SERVERADDRPORT = "localhost:8080"
	_STOREINTERVAL  = 300 * time.Second
	_STOREFILE      = "metrics-db.json"
	_RESTORE        = true

	_STOREINTERVALENV = "STORE_INTERVAL"
)

var (
	ErrParseDucation = errors.New("failed to parse duration")
	ErrUnknownEnvVar = errors.New("unknown env variable")
)

type Config struct {
	Address          string `env:"ADDRESS"`
	StoreInterval    time.Duration
	StoreFile        string `env:"STORE_FILE"`
	RestoreSavedData bool   `env:"RESTORE"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	Key              string `env:"KEY"`
}

func parseDuration(value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)

	if err != nil {
		if value, err := strconv.Atoi(value); err == nil {
			return time.Duration(value) * time.Second, nil
		}

		return 0, ErrParseDucation
	}

	return duration, nil
}

func parseDurationENV(p *time.Duration, envkey string) error {
	value, ok := os.LookupEnv(envkey)
	if !ok {
		return nil
	}

	valueDur, err := parseDuration(value)
	if err != nil {
		return err
	}

	*p = valueDur

	return nil

}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", _SERVERADDRPORT, "server address and port")
	flag.StringVar(&cfg.StoreFile, "f", _STOREFILE, "server db store file")
	flag.StringVar(&cfg.Key, "k", "", "key for digital sign")
	flag.BoolVar(&cfg.RestoreSavedData, "r", cfg.RestoreSavedData, "server restore db from file on start?")
	flag.Func("i", "server store interval", func(flagValue string) error {
		valueDur, err := parseDuration(flagValue)
		if err != nil {
			return err
		}
		cfg.StoreInterval = valueDur
		return nil
	})
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database connection string")

	flag.Parse()

	return nil
}

func getEnv(cfg *Config) error {

	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	if err = parseDurationENV(&cfg.StoreInterval, _STOREINTERVALENV); err != nil {
		return err
	}

	return nil

}

func Read() (Config, error) {
	var err error

	var cfg = Config{
		Address:          _SERVERADDRPORT,
		StoreInterval:    _STOREINTERVAL,
		StoreFile:        _STOREFILE,
		RestoreSavedData: _RESTORE,
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
