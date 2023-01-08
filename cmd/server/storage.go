package main

import "errors"

type MemStorage struct {
	gMetrics map[string]GaugeMetric
	cMetrics map[string]CounterMetric
}

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

var storage = createDB()

func createDB() *MemStorage {
	var d MemStorage
	d.gMetrics = make(map[string]GaugeMetric)
	d.cMetrics = make(map[string]CounterMetric)
	return &d
}

func updateGMetric(g GaugeMetric, s *MemStorage) {
	s.gMetrics[g.Name] = g
}

func updateCMetric(c CounterMetric, s *MemStorage) {
	s.cMetrics[c.Name] = c
}

func GetGMetric(mname string, s *MemStorage) (GaugeMetric, error) {
	if metric, ok := s.gMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func GetCMetric(mname string, s *MemStorage) (CounterMetric, error) {
	if metric, ok := s.cMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}
