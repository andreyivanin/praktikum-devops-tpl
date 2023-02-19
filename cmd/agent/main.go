package main

import (
	"devops-tpl/internal/agent"
	"fmt"
	"time"
)

func main() {
	cfg := agent.GetEnvConfig()
	requestTicker := time.NewTicker(cfg.PollInterval)
	sendTicker := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-requestTicker.C:
			// runtime.ReadMemStats(&Rtm)
			agent.PollMetrics()

			fmt.Println("Metric update", " - ", time.Now())
		case <-sendTicker.C:
			agent.SendMetricsJSON()

			// GMetricObjects := agent.GMetricGeneratorNew()
			// for _, object := range GMetricObjects {
			// 	go object.SendMetricJSON()
			// }
			// CMetricObjects := agent.CMetricGenerator(agent.PollCounter)
			// for _, object := range CMetricObjects {
			// 	go object.SendMetric()
			// }
		}
	}

}
