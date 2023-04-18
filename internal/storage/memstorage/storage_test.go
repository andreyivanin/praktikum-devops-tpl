package memstorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateMetric(t *testing.T) {

	type fields struct {
		name   string
		metric Metric
	}

	tests := []struct {
		name   string
		metric fields
		want   MemStorage
	}{
		{
			name: "update gauge metric",
			metric: fields{
				name:   "Alloc",
				metric: GaugeMetric{MType: "gauge", Value: 1223113},
			},
			want: MemStorage{
				Metrics: Metrics{"Alloc": GaugeMetric{MType: "gauge", Value: 1223113}},
				Mu:      new(sync.Mutex),
			},
		},
		{
			name: "update counter metric",
			metric: fields{
				name:   "RandomValue",
				metric: CounterMetric{MType: "counter", Delta: 67},
			},
			want: MemStorage{
				Metrics: Metrics{"RandomValue": CounterMetric{MType: "counter", Delta: 134}},
				Mu:      new(sync.Mutex),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var DB = MemStorage{
				Metrics: make(map[string]Metric),
				Mu:      new(sync.Mutex),
			}
			DB.UpdateMetric(tt.metric.name, tt.metric.metric)
			DB.UpdateMetric(tt.metric.name, tt.metric.metric)
			assert.Equal(t, tt.want, DB)
		})
	}
}
