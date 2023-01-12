package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	POLLINTERVAL   = 2
	REPORTINTERVAL = 10
)

func main() {
	var rtm runtime.MemStats
	var pollCounter int
	requestTicker := time.NewTicker(POLLINTERVAL * time.Second)
	sendTicker := time.NewTicker(REPORTINTERVAL * time.Second)

	for {
		select {
		case <-requestTicker.C:
			runtime.ReadMemStats(&rtm)
			fmt.Println("Metric update", " - ", time.Now())
			pollCounter++
		case <-sendTicker.C:
			GMetricObjects := GMetricGenerator(rtm)
			for _, object := range GMetricObjects {
				go object.SendMetric()
			}
			CMetricObjects := CMetricGenerator(pollCounter)
			for _, object := range CMetricObjects {
				go object.SendMetric()
			}
		}
	}
}
