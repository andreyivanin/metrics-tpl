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

	_POLLINTERVALENV   = "POLL_INTERVAL"
	_REPORTINTERVALENV = "REPORT_INTERVAL"

	_MAXSENDATTEMPTS = 3
)

var (
	ErrParseDucation = errors.New("failed to parse duration")
	ErrUnknownEnvVar = errors.New("unknown env variable")
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string `env:"KEY"`
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
	flag.StringVar(&cfg.Key, "k", "", "key for digital sign")

	flag.Func("r", "agent report interval", func(flagValue string) error {
		valueDur, err := parseDuration(flagValue)
		if err != nil {
			return err
		}
		cfg.ReportInterval = valueDur
		return nil
	})

	flag.Func("p", "agent poll interval", func(flagValue string) error {
		valueDur, err := parseDuration(flagValue)
		if err != nil {
			return err
		}
		cfg.PollInterval = valueDur
		return nil
	})

	flag.Parse()

	return nil
}

func getEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	if err = parseDurationENV(&cfg.PollInterval, _POLLINTERVALENV); err != nil {
		return err
	}

	if err = parseDurationENV(&cfg.ReportInterval, _REPORTINTERVALENV); err != nil {
		return err
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
