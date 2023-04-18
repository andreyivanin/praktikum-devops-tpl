package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"devops-tpl/internal/storage/memstorage"
)

func (h *Handler) MetricUpdate(w http.ResponseWriter, r *http.Request) {
	var metric memstorage.Metric
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	mvalue := chi.URLParam(r, "mvalue")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch mtype {
	case "gauge":
		mvalueconv, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			return
		}

		metric = memstorage.GaugeMetric{
			MType: "gauge",
			Value: mvalueconv,
		}

	case "counter":
		mvalueconv, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			return
		}
		metric = memstorage.CounterMetric{
			MType: "counter",
			Delta: mvalueconv,
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Bad metric type"))
		return
	}

	h.Storage.UpdateMetric(mname, metric)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The metric " + mname + " was updated"))

}

func (h *Handler) MetricGet(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	metric, err := h.Storage.GetMetric(mname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The metric isn't found"))
		return
	}
	w.WriteHeader(http.StatusOK)

	switch mtype {
	case "gauge":
		metric := metric.(memstorage.GaugeMetric)
		valuestring := fmt.Sprintf("%.9g", metric.Value)
		w.Write([]byte(valuestring))

	case "counter":
		metric := metric.(memstorage.CounterMetric)
		valuestring := strconv.Itoa(int(metric.Delta))
		w.Write([]byte(valuestring))

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Bad metric type"))
	}
}

func (h *Handler) MetricSummary(w http.ResponseWriter, r *http.Request) {
	metrics := h.Storage.GetAllMetrics()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for name, metric := range metrics {
		valuestring := fmt.Sprintf("%v", metric)
		w.Write([]byte(name + ": " + valuestring + "\n"))
	}
}
