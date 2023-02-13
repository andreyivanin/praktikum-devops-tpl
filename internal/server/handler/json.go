package handler

import (
	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/memstorage"
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

func MetricJSON(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		log.Println(err)
	}

	switch jsonmetric.MType {
	case "gauge":
		gmetric := memstorage.GaugeMetric{
			Name:  jsonmetric.ID,
			Value: *jsonmetric.Value,
		}
		s.UpdateGMetric(gmetric)

	case "counter":
		cmetric := memstorage.CounterMetric{
			Name:  jsonmetric.ID,
			Value: *jsonmetric.Delta,
		}
		s.UpdateCMetric(cmetric)

		updatedMetric, err := s.GetCMetric(jsonmetric.ID)
		if err != nil {
			log.Panicln(err)
		}
		jsonmetric.Delta = &updatedMetric.Value
	}

	metricsJSON, err := json.Marshal(jsonmetric)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}

func MetricSummaryJSON(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	MetricOK := true

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		log.Println(err)
	}

	switch jsonmetric.MType {
	case "gauge":
		if gmetric, err := s.GetGMetric(jsonmetric.ID); err != nil {
			log.Println(err)
			MetricOK = false
		} else {
			jsonmetric.Value = &gmetric.Value
		}
	case "counter":
		if cmetric, err := s.GetCMetric(jsonmetric.ID); err != nil {
			log.Println(err)
			MetricOK = false
		} else {
			jsonmetric.Delta = &cmetric.Value
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
