package main

import (
	"fmt"
	"log"
	"metrics-tpl/internal/agent"
	"time"
)

func main() {
	cfg, err := agent.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	monitor := agent.NewMonitor(cfg)

	for {
		select {
		case <-monitor.UpdateTicker.C:
			monitor.UpdateMetrics()
			fmt.Println("Metrics update", " - ", time.Now())
		case <-monitor.SendTicker.C:
			monitor.SendMetrics()
			fmt.Println("Metrics send", " - ", time.Now())
		}
	}
}
