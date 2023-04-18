package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devops-tpl/internal/agent"
)

func main() {
	cfg, err := agent.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	monitor := agent.NewMonitor(cfg)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-monitor.UpdateTicker.C:
			monitor.UpdateMetrics()
			fmt.Println("Metrics update", " - ", time.Now())
		case <-monitor.SendTicker.C:
			monitor.SendMetricsJSON()
			fmt.Println("Metrics send", " - ", time.Now())
		case sig := <-termSignal:
			monitor.UpdateTicker.Stop()
			monitor.SendTicker.Stop()
			log.Println("Finished, reason:", sig.String())
			os.Exit(0)
		}
	}
}
