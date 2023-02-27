package main

import (
	"devops-tpl/internal/agent"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := agent.GetConfig()
	requestTicker := time.NewTicker(cfg.PollInterval)
	sendTicker := time.NewTicker(cfg.ReportInterval)
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-requestTicker.C:
			agent.PollMetrics()
			fmt.Println("Metric update", " - ", time.Now())
		case <-sendTicker.C:
			agent.SendMetricsJSON(cfg)
		case sig := <-termSignal:
			log.Panicln("Finished, reason:", sig.String())
		}
	}
}
