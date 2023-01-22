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

type MetricsSlice []*Metrics

func (ms *MetricsSlice) UnmarshalJSON(data []byte) error {
	if data[0] != 0x5b { //check first "[" letter and make slice in string
		data = append([]byte{0x5b}, data...)
		data = append(data, 0x5d)
	}
	type MetricAlias struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}

	var MetricAliasSlice []MetricAlias
	if err := json.Unmarshal(data, &MetricAliasSlice); err != nil {
		log.Println(err)
	}

	for _, metric := range MetricAliasSlice {
		*ms = append(*ms, &Metrics{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: metric.Delta,
			Value: metric.Value,
		})
	}
	return nil
}

func MetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// jsonmetrics := []Metrics{}
	jsonmetrics := MetricsSlice{}

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
			storage.UpdateGMetric(gmetric, storage.DB)

		case "counter":
			cmetric := storage.CounterMetric{
				Name:  jsonmetric.ID,
				Value: *jsonmetric.Delta,
			}
			storage.UpdateCMetric(cmetric, storage.DB)
		}

	}

	// jsonmetrics = []Metrics{}
	// jsonmetrics = MetricsSlice{}
	// metrics := storage.GetMetricSummary()

	// for _, gMetric := range metrics.GMetrics {
	// 	name := gMetric.Name
	// 	value := gMetric.Value
	// 	jsonMetric := Metrics{
	// 		ID:    name,
	// 		MType: "gauge",
	// 		Value: &value,
	// 	}
	// 	jsonmetrics = append(jsonmetrics, &jsonMetric)

	// }

	// for _, cMetric := range metrics.CMetrics {
	// 	name := cMetric.Name
	// 	value := cMetric.Value
	// 	jsonMetric := Metrics{
	// 		ID:    name,
	// 		MType: "counter",
	// 		Delta: &value,
	// 	}
	// 	jsonmetrics = append(jsonmetrics, &jsonMetric)
	// }

	metricsJSON, err := json.Marshal(jsonmetrics)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}

func MetricSummaryJSONHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	MetricOK := true

	jsonmetrics := MetricsSlice{}

	if err := json.Unmarshal(b, &jsonmetrics); err != nil {
		log.Println(err)
	}

	for _, jsonmetric := range jsonmetrics {
		switch jsonmetric.MType {
		case "gauge":
			if gmetric, err := storage.GetGMetric(jsonmetric.ID); err != nil {
				log.Println(err)
				MetricOK = MetricOK && false
			} else {
				jsonmetric.Value = &gmetric.Value
			}
		case "counter":
			if cmetric, err := storage.GetCMetric(jsonmetric.ID); err != nil {
				log.Println(err)
				MetricOK = MetricOK && false
			} else {
				jsonmetric.Delta = &cmetric.Value
			}
		default:
			log.Println("wrong metric type")
			MetricOK = MetricOK && false
		}
	}

	if MetricOK {
		metricsJSON, err := json.Marshal(jsonmetrics)

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
