package agent

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

type GaugeMetric float64
type CounterMetric int64

type Monitor struct {
	cfg          Config
	values       runtime.MemStats
	pollCounter  int
	UpdateTicker *time.Ticker
	SendTicker   *time.Ticker
	Metrics      []Metric
}

func NewMonitor(cfg Config) Monitor {
	return Monitor{
		cfg:          cfg,
		UpdateTicker: time.NewTicker(cfg.PollInterval),
		SendTicker:   time.NewTicker(cfg.ReportInterval),
	}
}

func (m *Monitor) UpdateMetrics() {
	runtime.ReadMemStats(&m.values)
	m.pollCounter++

	Metrics := make(map[string]interface{}, 29)
	Metrics["Alloc"] = GaugeMetric(m.values.Alloc)
	Metrics["BuckHashSys"] = GaugeMetric(m.values.BuckHashSys)
	Metrics["Frees"] = GaugeMetric(m.values.Frees)
	Metrics["GCCPUFraction"] = GaugeMetric(m.values.GCCPUFraction)
	Metrics["GCSys"] = GaugeMetric(m.values.GCSys)
	Metrics["HeapAlloc"] = GaugeMetric(m.values.HeapAlloc)
	Metrics["HeapIdle"] = GaugeMetric(m.values.HeapIdle)
	Metrics["HeapInuse"] = GaugeMetric(m.values.HeapInuse)
	Metrics["HeapObjects"] = GaugeMetric(m.values.HeapObjects)
	Metrics["HeapReleased"] = GaugeMetric(m.values.HeapReleased)
	Metrics["HeapSys"] = GaugeMetric(m.values.HeapSys)
	Metrics["LastGC"] = GaugeMetric(m.values.LastGC)
	Metrics["Lookups"] = GaugeMetric(m.values.Lookups)
	Metrics["MCacheInuse"] = GaugeMetric(m.values.MCacheInuse)
	Metrics["MCacheSys"] = GaugeMetric(m.values.MCacheSys)
	Metrics["MSpanInuse"] = GaugeMetric(m.values.MSpanInuse)
	Metrics["MSpanSys"] = GaugeMetric(m.values.MSpanSys)
	Metrics["Mallocs"] = GaugeMetric(m.values.Mallocs)
	Metrics["NextGC"] = GaugeMetric(m.values.NextGC)
	Metrics["NumForcedGC"] = GaugeMetric(m.values.NumForcedGC)
	Metrics["NumGC"] = GaugeMetric(m.values.NumGC)
	Metrics["OtherSys"] = GaugeMetric(m.values.OtherSys)
	Metrics["PauseTotalNs"] = GaugeMetric(m.values.PauseTotalNs)
	Metrics["StackInuse"] = GaugeMetric(m.values.StackInuse)
	Metrics["StackSys"] = GaugeMetric(m.values.StackSys)
	Metrics["Sys"] = GaugeMetric(m.values.Sys)
	Metrics["TotalAlloc"] = GaugeMetric(m.values.TotalAlloc)
	Metrics["RandomValue"] = GaugeMetric(rand.Intn(100))
	Metrics["PollCount"] = CounterMetric(m.pollCounter)

	for name, value := range Metrics {
		switch value.(type) {
		case GaugeMetric:
			m.Metrics = append(m.Metrics, Metric{
				name:  name,
				mtype: "gauge",
				value: value.(GaugeMetric),
			})
		case CounterMetric:
			m.Metrics = append(m.Metrics, Metric{
				name:  name,
				mtype: "counter",
				delta: value.(CounterMetric),
			})
		}
	}
}

func (m *Monitor) SendMetrics() {
	client := http.Client{}

	for _, metric := range m.Metrics {
		url := metric.CreateURL(m.cfg)
		log.Println(url)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)

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
}

type Metric struct {
	name  string
	mtype string
	value GaugeMetric
	delta CounterMetric
}

func (m *Metric) CreateURL(cfg Config) string {
	var u url.URL
	var valuestring string
	u.Scheme = PROTOCOL
	u.Host = cfg.Address

	switch m.mtype {
	case "gauge":
		valuestring = fmt.Sprintf("%.f", m.value)
	case "counter":
		valuestring = strconv.Itoa(int(m.delta))
	}

	url := u.JoinPath("update", m.mtype, m.name, valuestring)
	return url.String()
}
