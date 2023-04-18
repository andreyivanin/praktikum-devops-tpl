package agent

import (
	"devops-tpl/internal/storage/memstorage"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMetricsJSON(t *testing.T) {
	type fields struct {
		url  string
		body interface{}
	}

	tests := []struct {
		name    string
		metrics []Metric
		want    fields
	}{
		{
			name: "good test#1: Gaugemetric",
			metrics: []Metric{
				Metric{name: "Alloc", mtype: "gauge", value: 150},
				// Metric{name: "PollCount", mtype: "counter", delta: 55},
			},
			want: fields{
				url:  "/update",
				body: memstorage.GaugeMetric{MType: "gauge", Value: 150},
			},
		},
		{
			name: "good test#2: Countermetric",
			metrics: []Metric{
				Metric{name: "PollCount", mtype: "counter", delta: 55},
			},
			want: fields{
				url:  "/update",
				body: memstorage.CounterMetric{MType: "counter", Delta: 55},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mon := Monitor{
				cfg:     Config{Address: "127.0.0.1:8080"},
				Metrics: tt.metrics,
			}

			l, err := net.Listen("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				path := req.URL.Path

				b, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
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

				assert.Equal(t, tt.want.url, path)
				assert.Equal(t, tt.want.body, metric)
			}))
			defer func() { ts.Close() }()

			ts.Listener.Close()
			ts.Listener = l
			ts.Start()

			mon.SendMetricsJSON()

		})
	}
}
