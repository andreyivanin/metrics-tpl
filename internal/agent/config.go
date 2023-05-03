package agent

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	PROTOCOL       = "http"
	SERVERADDRPORT = "localhost:8080"
	POLLINTERVAL   = 2
	REPORTINTERVAL = 10
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func GetFlagConfig(cfg *Config) {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address and port")
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "agent report interval")
	flag.IntVar(&cfg.PollInterval, "p", cfg.PollInterval, "agent poll interval")
	flag.Parse()
}

func GetEnvConfig(cfg *Config) {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() (Config, error) {
	var cfg = Config{
		Address:        SERVERADDRPORT,
		PollInterval:   POLLINTERVAL,
		ReportInterval: REPORTINTERVAL,
	}

	GetFlagConfig(&cfg)
	GetEnvConfig(&cfg)

	return cfg, nil
}
