package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"devops-tpl/internal/storage/memstorage"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (h *Handler) MetricJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		log.Println(err)
	}

	var metric memstorage.Metric

	switch jsonmetric.MType {
	case "gauge":
		metric = memstorage.GaugeMetric{
			MType: "gauge",
			Value: *jsonmetric.Value,
		}

	case "counter":
		metric = memstorage.CounterMetric{
			MType: "counter",
			Delta: *jsonmetric.Delta,
		}

	}

	updatedMetric, err := h.Storage.UpdateMetric(jsonmetric.ID, metric)
	if err != nil {
		log.Println(err)
		return
	}

	switch updatedMetric := updatedMetric.(type) {
	case memstorage.GaugeMetric:
		jsonmetric.Value = &updatedMetric.Value

	case memstorage.CounterMetric:
		jsonmetric.Delta = &updatedMetric.Delta
	}

	metricsJSON, err := json.Marshal(jsonmetric)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}

func (h *Handler) MetricSummaryJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	MetricOK := true

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		log.Println(err)
	}

	switch jsonmetric.MType {
	case "gauge":
		if metric, err := h.Storage.GetMetric(jsonmetric.ID); err != nil {
			log.Println(err)
			MetricOK = false
		} else {
			metric := metric.(memstorage.GaugeMetric)
			jsonmetric.Value = &metric.Value
		}
	case "counter":
		if metric, err := h.Storage.GetMetric(jsonmetric.ID); err != nil {
			log.Println(err)
			MetricOK = false
		} else {
			metric := metric.(memstorage.CounterMetric)
			jsonmetric.Delta = &metric.Delta
		}
	default:
		log.Println("wrong metric type")
		MetricOK = false
	}

	if MetricOK {
		metricsJSON, err := json.Marshal(jsonmetric)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(metricsJSON)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("One or several metrics weren't found"))
	}
}
