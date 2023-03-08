package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"devops-tpl/internal/storage"
	"devops-tpl/internal/storage/memstorage"

	"github.com/go-chi/chi/v5"
)

func MetricUpdate(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	mvalue := chi.URLParam(r, "mvalue")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch mtype {
	case "gauge":
		floatvalue, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		gmetric := memstorage.GaugeMetric{
			Name:  mname,
			Value: floatvalue,
		}
		s.UpdateGMetric(gmetric)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + gmetric.Name + " was updated"))

	case "counter":
		intvalue, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		cmetric := memstorage.CounterMetric{
			Name:  mname,
			Value: intvalue,
		}
		s.UpdateCMetric(cmetric)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + cmetric.Name + " was updated"))

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Bad metric type"))
	}

}

func MetricGet(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch mtype {
	case "gauge":
		metric, err := s.GetGMetric(mname)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("The metric isn't found"))
			break
		}
		w.WriteHeader(http.StatusOK)
		valuestring := fmt.Sprintf("%.9g", metric.Value)
		w.Write([]byte(valuestring))

	case "counter":
		metric, err := s.GetCMetric(mname)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("The metric isn't found"))
			break
		}
		w.WriteHeader(http.StatusOK)
		valuestring := strconv.Itoa(int(metric.Value))
		w.Write([]byte(valuestring))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Bad metric type"))
	}
}

func MetricSummary(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	metrics := s.GetStorage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, metric := range metrics.GMetrics {
		valuestring := fmt.Sprintf("%.f", metric.Value)
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}

	for _, metric := range metrics.CMetrics {
		valuestring := strconv.Itoa(int(metric.Value))
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}
}
