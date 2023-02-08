package storage

import (
	"devops-tpl/internal/storage/memory"
)

type Storage interface {
	UpdateGMetric(memory.GaugeMetric)
	UpdateCMetric(memory.CounterMetric)
	GetGMetric(string) (memory.GaugeMetric, error)
	GetCMetric(string) (*memory.CounterMetric, error)
	GetMetricSummary()
}
