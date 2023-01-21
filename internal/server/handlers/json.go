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
			storage.UpdateGMetric(gmetric, storage.DB)

		case "counter":
			cmetric := storage.CounterMetric{
				Name:  jsonmetric.ID,
				Value: *jsonmetric.Delta,
			}
			storage.UpdateCMetric(cmetric, storage.DB)
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

	if string(b)[0:1] == "[" {
		jsonmetrics := []*Metrics{}
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
	} else {
		jsonmetric := Metrics{}
		if err := json.Unmarshal(b, &jsonmetric); err != nil {
			log.Println(err)
		}
		switch jsonmetric.MType {
		case "gauge":
			if gmetric, err := storage.GetGMetric(jsonmetric.ID); err != nil {
				log.Println(err)
				MetricOK = false
			} else {
				jsonmetric.Value = &gmetric.Value
			}
		case "counter":
			if cmetric, err := storage.GetCMetric(jsonmetric.ID); err != nil {
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
}

// jsonmetrics := []*Metrics{}

// if err := json.Unmarshal(b, &jsonmetrics); err != nil {
// 	log.Println(err)
// }

// 	for _, jsonmetric := range jsonmetrics {
// 		switch jsonmetric.MType {
// 		case "gauge":
// 			if gmetric, err := storage.GetGMetric(jsonmetric.ID); err == nil {
// 				jsonmetric.Value = &gmetric.Value
// 			}
// 		case "counter":
// 			if cmetric, err := storage.GetCMetric(jsonmetric.ID); err == nil {
// 				jsonmetric.Delta = &cmetric.Value
// 			}
// 		}
// 	}

// 	metricsJSON, err := json.Marshal(jsonmetrics)

// 	if err != nil {
// 		panic(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(metricsJSON)
// }

// type Slice []byte

// type Result struct {
// 	Key   string
// 	Value []int
// }

// func (s *Slice) UnmarshalJSON(data []byte) error {
// 	var obj map[string]json.RawMessage
// 	if err := json.Unmarshal(data, &obj); err != nil {
// 		return err
// 	}

// 	for key, raw := range obj {
// 		r := Result{Key: key}
// 		if raw[0] == '[' {
// 			if err := json.Unmarshal(raw, &r.Value); err != nil {
// 				return err
// 			}
// 		} else {
// 			var i int
//             if err := json.Unmarshal(raw, &i); err != nil {
//                 return err
//             }
//             r.Value = append(r.Value, i)
//         }
// 		*s = append(*s, r)

// 	}
// 	return nil

// }

// func (s *Slice) UnmarshalJSON(data []byte) error {
// 	if s[0] == "[" {
// 		jsonmetrics := []*Metrics{}
// 		if err := json.Unmarshal(s, jsonmetrics); err != nil {
// 			return nil
// 		}
// 	} else {
// 		jsonmetrics = *Metrics{}
// 	}
// }
