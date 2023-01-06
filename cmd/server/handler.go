package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func metricHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	fields := strings.Split(url, "/")
	if len(fields) == 5 {
		switch fields[2] {
		case "gauge":
			floatvalue, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad metric value"))
				break
			}

			gmetric := GaugeMetric{
				Name:  fields[3],
				Value: floatvalue,
			}
			updateGMetric(gmetric, storage)
			fmt.Print(storage)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("The metric " + gmetric.Name + " was updated"))

		case "counter":
			intvalue, err := strconv.ParseInt(fields[4], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad metric value"))
				break
			}

			cmetric := CounterMetric{
				Name:  fields[3],
				Value: intvalue,
			}
			updateCMetric(cmetric, storage)
			fmt.Print(storage)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("The metric " + cmetric.Name + " was updated"))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown metric type"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Bad URL"))
	}

}
