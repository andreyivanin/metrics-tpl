package main

import (
	"fmt"
	"metrics-tpl/internal/agent"
	"time"
)

func main() {
	monitor := agent.NewMonitor()

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
