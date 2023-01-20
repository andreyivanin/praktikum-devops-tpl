package handlers

import (
	"devops-tpl/internal/storage"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func MetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	jsonmetrics := []Metrics{}
	if err := json.Unmarshal(b, &jsonmetrics); err != nil {
		log.Println(err)
	}

	for _, jsonmetric := range jsonmetrics {
		switch jsonmetric.MType {
		case "gauge":
			gmetric := storage.GaugeMetric{
				Name:  jsonmetric.ID,
				Value: *jsonmetric.Value,
			}
			storage.UpdateGMetric(gmetric)

		case "counter":
			cmetric := storage.CounterMetric{
				Name:  jsonmetric.ID,
				Value: *jsonmetric.Delta,
			}
			storage.UpdateCMetric(cmetric)
		}

	}

	jsonmetrics = []Metrics{}
	metrics := storage.GetMetricSummary()

	for _, gMetric := range metrics.GMetrics {
		name := gMetric.Name
		value := gMetric.Value
		jsonMetric := Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		jsonmetrics = append(jsonmetrics, jsonMetric)

	}

	for _, cMetric := range metrics.CMetrics {
		name := cMetric.Name
		value := cMetric.Value
		jsonMetric := Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}
		jsonmetrics = append(jsonmetrics, jsonMetric)
	}

	metricsJSON, err := json.Marshal(jsonmetrics)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}

func MetricSummaryJSONHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	jsonmetrics := []*Metrics{}
	if err := json.Unmarshal(b, &jsonmetrics); err != nil {
		log.Println(err)
	}

	for _, jsonmetric := range jsonmetrics {
		switch jsonmetric.MType {
		case "gauge":
			if gmetric, err := storage.GetGMetric(jsonmetric.ID); err == nil {
				jsonmetric.Value = &gmetric.Value
			}
		case "counter":
			if cmetric, err := storage.GetCMetric(jsonmetric.ID); err == nil {
				jsonmetric.Delta = &cmetric.Value
			}
		}
	}

	metricsJSON, err := json.Marshal(jsonmetrics)

	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}
