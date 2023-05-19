package config

import (
	"errors"
	"flag"
	"fmt"
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
}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", _SERVERADDRPORT, "server address and port")
	flag.StringVar(&cfg.StoreFile, "f", _STOREFILE, "server db store file")
	flag.BoolVar(&cfg.RestoreSavedData, "r", cfg.RestoreSavedData, "server restore db from file on start?")
	// flag.DurationVar(&cfg.StoreInterval, "i", _STOREINTERVAL, "server store interval")
	flag.Func("i", "server store interval", func(flagValue string) error {

		if len(flagValue) > 0 {
			if _, err := strconv.Atoi(flagValue); err == nil {
				flagValue += "s"
			}
		}

		d, err := time.ParseDuration(flagValue)
		if err != nil {
			return errors.New("failed to parse duration")
		}

		cfg.StoreInterval = d

		return nil

		// if err != nil {
		// 	if value, err := strconv.Atoi(flagValue); err == nil {
		// 		return time.Duration(flagValue) * time.Second, nil
		// 	}

		// 	return 0, ErrParseDucation
		// }

	})

	flag.Parse()

	return nil
}

func parseDuration(key string) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, ErrParseDucation
	}
	duration, err := time.ParseDuration(value)

	if err != nil {
		if value, err := strconv.Atoi(value); err == nil {
			return time.Duration(value) * time.Second, nil
		}

		return 0, ErrParseDucation
	}

	return duration, nil
}

func getEnv(cfg *Config) error {

	// parseDuration := func(v string) (interface{}, error) {
	// 	if len(v) > 0 {
	// 		if _, err := strconv.Atoi(v); err == nil {
	// 			v += "s"
	// 		}
	// 	}

	// 	d, err := time.ParseDuration(v)
	// 	if err != nil {
	// 		return nil, errors.New("failed to parse duration")
	// 	}

	// 	return d, nil
	// }

	// opt := env.Options{
	// 	FuncMap: map[reflect.Type]env.ParserFunc{
	// 		reflect.TypeOf(cfg.StoreInterval): parseDuration,
	// 	},
	// }

	// err := env.ParseWithOptions(cfg, opt)
	// if err != nil {
	// 	return err
	// }

	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	durationEnvs := [...]string{"STORE_INTERVAL"}

	for _, env := range durationEnvs {

		envDuration, err := parseDuration(env)
		if err != nil {
			return err
		}

		switch env {
		case "STORE_INTERVAL":
			cfg.StoreInterval = envDuration
		default:
			return ErrUnknownEnvVar
		}
	}

	return nil
}

func Read() (Config, error) {
	for i, v := range os.Args[1:] {
		fmt.Println(i+1, v)
	}

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
