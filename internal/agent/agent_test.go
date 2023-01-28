package agent

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeMetric_SendMetric(t *testing.T) {
	type fields struct {
		Name  string
		Value float64
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good test#1",
			fields: fields{
				Name:  "Alloc",
				Value: 100,
			},
			want: "/update/gauge/Alloc/100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GaugeMetric{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}

			l, err := net.Listen("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				path := req.URL.Path
				assert.Equal(t, path, tt.want)
			}))
			defer func() { ts.Close() }()

			ts.Listener.Close()
			ts.Listener = l
			ts.Start()
			g.SendMetric()

		})
	}
}
func TestCounterMetric_SendMetric(t *testing.T) {
	type fields struct {
		Name  string
		Value int64
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good test#1",
			fields: fields{
				Name:  "RandomValue",
				Value: 67,
			},
			want: "/update/counter/RandomValue/67",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CounterMetric{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}

			l, err := net.Listen("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				path := req.URL.Path
				assert.Equal(t, path, tt.want)
			}))
			defer func() { ts.Close() }()

			ts.Listener.Close()
			ts.Listener = l
			ts.Start()
			c.SendMetric()

		})
	}
}
