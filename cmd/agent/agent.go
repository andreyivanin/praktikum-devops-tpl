package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
)

const (
	PROTOCOL       = "http"
	SERVERADDRPORT = "127.0.0.1:8080"
)

type Monitor struct {
	Alloc,
	BuckHashSys,
	Frees GaugeMetric
}

type GaugeMetric struct {
	Name  string
	Value float64
}

func (g GaugeMetric) SendMetric() {
	client := http.Client{}
	url := g.CreateURL()
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	if response != nil {
		fmt.Println("Status code", response.Status)

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(body))
	}
}

func (g GaugeMetric) CreateURL() string {
	var u url.URL
	valuestring := fmt.Sprintf("%.0f", g.Value)
	u.Scheme = PROTOCOL
	u.Host = SERVERADDRPORT
	url := u.JoinPath("update", "gauge", g.Name, valuestring)
	return url.String()

}

type CounterMetric struct {
	Name  string
	Value int64
}

func (c CounterMetric) SendMetric() {
	client := http.Client{}
	url := c.CreateURL()
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	if response != nil {
		fmt.Println("Status code", response.Status)

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(body))
	}
}

func (c CounterMetric) CreateURL() string {
	var u url.URL
	valuestring := strconv.Itoa(int(c.Value))
	u.Scheme = PROTOCOL
	u.Host = SERVERADDRPORT
	url := u.JoinPath("update", "counter", c.Name, valuestring)
	return url.String()
}

func CreateGM(name string, value float64) *GaugeMetric {
	return &GaugeMetric{
		Name:  name,
		Value: value,
	}
}

func CreateCM(name string, value int64) *CounterMetric {
	return &CounterMetric{
		Name:  name,
		Value: value,
	}
}

func GMetricGenerator(values runtime.MemStats) []*GaugeMetric {
	gMetrics := make(map[string]float64)
	gMetrics["Alloc"] = float64(values.Alloc)
	gMetrics["BuckHashSys"] = float64(values.BuckHashSys)
	gMetrics["Frees"] = float64(values.Frees)
	gMetrics["GCCPUFraction"] = float64(values.GCCPUFraction)
	gMetrics["GCSys"] = float64(values.GCSys)
	gMetrics["HeapAlloc"] = float64(values.HeapAlloc)
	gMetrics["HeapIdle"] = float64(values.HeapIdle)
	gMetrics["HeapInuse"] = float64(values.HeapInuse)
	gMetrics["HeapObjects"] = float64(values.HeapObjects)
	gMetrics["HeapReleased"] = float64(values.HeapReleased)
	gMetrics["HeapSys"] = float64(values.HeapSys)
	gMetrics["LastGC"] = float64(values.LastGC)
	gMetrics["Lookups"] = float64(values.Lookups)
	gMetrics["MCacheInuse"] = float64(values.MCacheInuse)
	gMetrics["MCacheSys"] = float64(values.MCacheSys)
	gMetrics["MSpanInuse"] = float64(values.MSpanInuse)
	gMetrics["MSpanSys"] = float64(values.MSpanSys)
	gMetrics["Mallocs"] = float64(values.Mallocs)
	gMetrics["NextGC"] = float64(values.NextGC)
	gMetrics["NumForcedGC"] = float64(values.NumForcedGC)
	gMetrics["NumGC"] = float64(values.NumGC)
	gMetrics["OtherSys"] = float64(values.OtherSys)
	gMetrics["PauseTotalNs"] = float64(values.PauseTotalNs)
	gMetrics["StackInuse"] = float64(values.StackInuse)
	gMetrics["StackSys"] = float64(values.StackSys)
	gMetrics["Sys"] = float64(values.Sys)
	gMetrics["TotalAlloc"] = float64(values.TotalAlloc)

	GMetricObjects := []*GaugeMetric{}
	for name, value := range gMetrics {
		object := CreateGM(name, value)
		GMetricObjects = append(GMetricObjects, object)
	}
	return GMetricObjects
}

func CMetricGenerator(pollCounter int) []*CounterMetric {
	cMetrics := make(map[string]int64)

	cMetrics["PollCount"] = int64(pollCounter)
	cMetrics["RandomValue"] = int64(rand.Intn(100))

	CMetricObjects := []*CounterMetric{}
	for name, value := range cMetrics {
		object := CreateCM(name, value)
		CMetricObjects = append(CMetricObjects, object)
	}
	return CMetricObjects

}
