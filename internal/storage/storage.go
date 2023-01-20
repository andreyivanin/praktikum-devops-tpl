package storage

import (
	"errors"
	"fmt"
)

type MemStorage struct {
	GMetrics map[string]GaugeMetric
	CMetrics map[string]*CounterMetric
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
	d.GMetrics = make(map[string]GaugeMetric)
	d.CMetrics = make(map[string]*CounterMetric)
	return &d
}

func UpdateGMetric(g GaugeMetric) {
	storage.GMetrics[g.Name] = g
}

func UpdateCMetric(c CounterMetric) {
	if metric, ok := storage.CMetrics[c.Name]; ok {
		metric.Value = metric.Value + c.Value
	} else {
		storage.CMetrics[c.Name] = &c
	}
	fmt.Println("ok")
}

func GetGMetric(mname string) (GaugeMetric, error) {
	if metric, ok := storage.GMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func GetCMetric(mname string) (*CounterMetric, error) {
	if metric, ok := storage.CMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func GetMetricSummary() *MemStorage {
	return storage
}
