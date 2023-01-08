package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", metricUpdateHandler)
			r.Get("/", metricUpdateHandler)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", metricGetHandler)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", metricSummaryHandler)
	})

	return r
}

func metricUpdateHandler(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	mvalue := chi.URLParam(r, "mvalue")
	//url := r.URL.Path
	// fields := strings.Split(url, "/")
	switch mtype {
	case "gauge":
		floatvalue, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		gmetric := GaugeMetric{
			Name:  mname,
			Value: floatvalue,
		}
		updateGMetric(gmetric, storage)
		fmt.Print(storage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + gmetric.Name + " was updated"))

	case "counter":
		intvalue, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			break
		}

		cmetric := CounterMetric{
			Name:  mname,
			Value: intvalue,
		}
		updateCMetric(cmetric, storage)
		fmt.Print(storage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + cmetric.Name + " was updated"))
	}
}

func metricGetHandler(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	switch mtype {
	case "gauge":
		metric, err := GetGMetric(mname, storage)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("The metric isn't found"))
			break
		}
		w.WriteHeader(http.StatusOK)
		valuestring := fmt.Sprintf("%.0f", metric.Value)
		w.Write([]byte(metric.Name + ": " + valuestring))

	case "counter":
		metric, err := GetCMetric(mname, storage)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("The metric isn't found"))
			break
		}
		w.WriteHeader(http.StatusOK)
		valuestring := strconv.Itoa(int(metric.Value))
		w.Write([]byte(metric.Name + ": " + valuestring))
	}
}

func metricSummaryHandler(w http.ResponseWriter, r *http.Request) {
	for _, metric := range storage.gMetrics {
		valuestring := fmt.Sprintf("%.0f", metric.Value)
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}

	for _, metric := range storage.cMetrics {
		valuestring := strconv.Itoa(int(metric.Value))
		w.Write([]byte(metric.Name + ": " + valuestring + "\n"))
	}
}
