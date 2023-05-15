package main

import (
	"log"
	"metrics-tpl/internal/agent"
)

func main() {
	cfg, err := agent.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	monitor := agent.NewMonitor(cfg)

	err = monitor.Run()
	if err != nil {
		log.Fatal(err)
	}
}
