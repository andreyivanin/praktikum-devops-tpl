package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"devops-tpl/internal/storage/memory"
)

func MetricUpdateHandler(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	mvalue := chi.URLParam(r, "mvalue")
	switch mtype {
	case "gauge":
		floatvalue, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		gmetric := memory.GaugeMetric{
			Name:  mname,
			Value: floatvalue,
		}
		memory.DB.UpdateGMetric(gmetric)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + gmetric.Name + " was updated"))

	case "counter":
		intvalue, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		cmetric := memory.CounterMetric{
			Name:  mname,
			Value: intvalue,
		}
		memory.DB.UpdateCMetric(cmetric)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + cmetric.Name + " was updated"))

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Bad metric type"))
	}

}

func MetricGetHandler(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	switch mtype {
	case "gauge":
		metric, err := memory.DB.GetGMetric(mname)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("The metric isn't found"))
			break
		}
		w.WriteHeader(http.StatusOK)
		valuestring := fmt.Sprintf("%.9g", metric.Value)
		w.Write([]byte(valuestring))

	case "counter":
		metric, err := memory.DB.GetCMetric(mname)
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

func MetricSummaryHandler(w http.ResponseWriter, r *http.Request) {
	metrics := memory.DB.GetMetricSummary()
	for _, metric := range metrics.GMetrics {
		valuestring := fmt.Sprintf("%.f", metric.Value)
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}

	for _, metric := range metrics.CMetrics {
		valuestring := strconv.Itoa(int(metric.Value))
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}
}
