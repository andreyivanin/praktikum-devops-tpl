package memstorage

import (
	"errors"
	"fmt"
	"sync"
)

// type Gauge float64

// func (g Gauge) MakePointer() *float64 {
// 	p := float64(g)
// 	return &p
// }

// type Counter int64

// func (c Counter) MakePointer() *int64 {
// 	p := int64(c)
// 	return &p
// }

type GaugeMetric struct {
	MType string
	Value float64
}

type CounterMetric struct {
	MType string
	Delta int64
}

type Metric interface{}

// func (m Metric) UnmarshalJSON(data []byte) error {
// 	fmt.Println(m)
// }

type Metrics map[string]Metric

func (m *Metrics) UnmarshalJSON(data []byte) error {
	type MetricsAlias Metrics
	MetricFile := &struct {
		*MetricsAlias
	}{MetricsAlias: (*MetricsAlias)(m)}
	fmt.Println(MetricFile)
	return nil
}

type MemStorage struct {
	Metrics map[string]Metric
	Mu      *sync.Mutex
}

func New() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]Metric),
		Mu:      new(sync.Mutex),
	}
}

func (s *MemStorage) UpdateMetric(name string, m Metric) (Metric, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	switch m.(type) {
	case GaugeMetric:
		s.Metrics[name] = m
	case CounterMetric:
		if existingMetric, ok := s.Metrics[name]; ok {
			updatedDelta := existingMetric.(CounterMetric).Delta + m.(CounterMetric).Delta
			s.Metrics[name] = CounterMetric{
				MType: "counter",
				Delta: updatedDelta,
			}
		} else {
			s.Metrics[name] = m
		}
	default:
		return nil, errors.New("the metric isn't found")
	}

	return s.Metrics[name], nil
}

func (s *MemStorage) GetMetric(mname string) (Metric, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if metric, ok := s.Metrics[mname]; ok {
		return metric, nil
	}

	return nil, errors.New("the metric isn't found")
}

func (s *MemStorage) GetAllMetrics() Metrics {
	return s.Metrics

}
