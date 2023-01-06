package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateGMetric(t *testing.T) {

	type gMetrics map[string]GaugeMetric
	type cMetrics map[string]CounterMetric

	tests := []struct {
		name    string
		gmetric GaugeMetric
		want    MemStorage
	}{
		{
			name:    "update gauge metric",
			gmetric: GaugeMetric{Name: "Alloc", Value: 1223113},
			want: MemStorage{
				gMetrics: gMetrics{
					"Alloc": GaugeMetric{Name: "Alloc", Value: 1223113},
				},
				cMetrics: cMetrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var storage = MemStorage{
				gMetrics: make(map[string]GaugeMetric),
				cMetrics: make(map[string]CounterMetric),
			}
			updateGMetric(tt.gmetric, &storage)
			updateGMetric(tt.gmetric, &storage)

			assert.Equal(t, storage, tt.want)
		})
	}
}

func Test_updateCMetric(t *testing.T) {

	type gMetrics map[string]GaugeMetric
	type cMetrics map[string]CounterMetric

	tests := []struct {
		name    string
		gmetric GaugeMetric
		cmetric CounterMetric
		want    MemStorage
	}{
		{
			name:    "update counter metric",
			cmetric: CounterMetric{Name: "RandomValue", Value: 67},
			want: MemStorage{
				gMetrics: gMetrics{},
				cMetrics: cMetrics{
					"RandomValue": CounterMetric{Name: "RandomValue", Value: 67},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var storage = MemStorage{
				gMetrics: make(map[string]GaugeMetric),
				cMetrics: make(map[string]CounterMetric),
			}
			updateCMetric(tt.cmetric, &storage)
			updateCMetric(tt.cmetric, &storage)

			assert.Equal(t, storage, tt.want)
		})
	}
}
