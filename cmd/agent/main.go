package main

import (
	"devops-tpl/internal/agent"
	"fmt"
	"time"
)

const (
	POLLINTERVAL   = 2
	REPORTINTERVAL = 10
)

func main() {
	requestTicker := time.NewTicker(POLLINTERVAL * time.Second)
	sendTicker := time.NewTicker(REPORTINTERVAL * time.Second)

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
