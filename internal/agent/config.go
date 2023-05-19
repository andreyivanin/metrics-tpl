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
	_POLLINTERVAL   = 2 * time.Second
	_REPORTINTERVAL = 10 * time.Second
)

var (
	ErrParseDucation = errors.New("failed to parse duration")
	ErrUnknownEnvVar = errors.New("unknown env variable")
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func getFlag(cfg *Config) error {
	flag.StringVar(&cfg.Address, "a", _SERVERADDRPORT, "server address and port")
	flag.DurationVar(&cfg.ReportInterval, "r", _REPORTINTERVAL, "agent report interval")
	flag.DurationVar(&cfg.PollInterval, "p", _POLLINTERVAL, "agent poll interval")

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
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	durationEnvs := [...]string{
		"REPORT_INTERVAL",
		"POLL_INTERVAL",
	}

	for _, env := range durationEnvs {

		envDuration, err := parseDuration(env)
		if err != nil {
			return err
		}

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
		PollInterval:   _POLLINTERVAL,
		ReportInterval: _REPORTINTERVAL,
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
