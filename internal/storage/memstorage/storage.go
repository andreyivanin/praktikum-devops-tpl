package memstorage

import (
	"errors"
	"sync"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type MemStorage struct {
	GMetrics map[string]GaugeMetric
	CMetrics map[string]*CounterMetric
	sync.Mutex
}

// var DB = New()

func New() *MemStorage {
	return &MemStorage{
		GMetrics: make(map[string]GaugeMetric),
		CMetrics: make(map[string]*CounterMetric),
	}
}

func (s *MemStorage) UpdateGMetric(g GaugeMetric) {
	s.Lock()
	defer s.Unlock()
	s.GMetrics[g.Name] = g
}

func (s *MemStorage) UpdateCMetric(c CounterMetric) {
	s.Lock()
	defer s.Unlock()
	if existingMetric, ok := s.CMetrics[c.Name]; ok {
		existingMetric.Value = existingMetric.Value + c.Value
	} else {
		s.CMetrics[c.Name] = &c
	}
}

func (s *MemStorage) GetGMetric(mname string) (GaugeMetric, error) {
	s.Lock()
	defer s.Unlock()
	if metric, ok := s.GMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func (s *MemStorage) GetCMetric(mname string) (*CounterMetric, error) {
	s.Lock()
	defer s.Unlock()
	if metric, ok := s.CMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func (s *MemStorage) GetStorage() *MemStorage {
	return s
}
