package agent

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	_PROTOCOL       = "http"
	_SERVERADDRPORT = "localhost:8080"
	_POLLINTERVAL   = 2
	_REPORTINTERVAL = 10
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address and port")

	ReportIntervalFlag := flag.Int("r", _REPORTINTERVAL, "agent report interval")
	PollIntervalFlag := flag.Int("p", _POLLINTERVAL, "agent poll interval")

	flag.Parse()

	cfg.ReportInterval = time.Duration(*ReportIntervalFlag) * time.Second
	cfg.PollInterval = time.Duration(*PollIntervalFlag) * time.Second

	return nil
}

func getEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	durationEnvs := [...]string{
		"REPORT_INTERVAL",
		"POLL_INTERVAL",
	}

	for _, env := range durationEnvs {
		envString, ok := os.LookupEnv(env)

		if !ok {
			continue
		}

		envInt, err := strconv.Atoi(envString)
		if err != nil {
			return err
		}

		envDuration := time.Duration(envInt) * time.Second

		switch env {
		case "REPORT_INTERVAL":
			cfg.ReportInterval = envDuration
		case "POLL_INTERVAL":
			cfg.PollInterval = envDuration
		default:
			return errors.New("unknown env variable")
		}
	}

	return nil
}

func GetConfig() (Config, error) {
	var err error

	var cfg = Config{
		Address:        _SERVERADDRPORT,
		PollInterval:   _POLLINTERVAL * time.Second,
		ReportInterval: _REPORTINTERVAL * time.Second,
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
