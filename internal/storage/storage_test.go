package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateGMetric(t *testing.T) {

	type gMetrics map[string]GaugeMetric
	type cMetrics map[string]*CounterMetric

	tests := []struct {
		name    string
		gmetric GaugeMetric
		want    MemStorage
	}{
		{
			name:    "update gauge metric",
			gmetric: GaugeMetric{Name: "Alloc", Value: 1223113},
			want: MemStorage{
				GMetrics: gMetrics{
					"Alloc": GaugeMetric{Name: "Alloc", Value: 1223113},
				},
				CMetrics: cMetrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateGMetric(tt.gmetric)
			UpdateGMetric(tt.gmetric)
			assert.Equal(t, tt.want, *storage)
		})
	}
}

func Test_updateCMetric(t *testing.T) {

	type gMetrics map[string]GaugeMetric
	type cMetrics map[string]*CounterMetric

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
				GMetrics: gMetrics{},
				CMetrics: cMetrics{
					"RandomValue": &CounterMetric{Name: "RandomValue", Value: 134},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateCMetric(tt.cmetric)
			UpdateCMetric(tt.cmetric)
			assert.Equal(t, tt.want, *storage)
		})
	}
}
