package storage

import (
	"devops-tpl/internal/storage/memstorage"
)

type Storage interface {
	UpdateMetric(string, memstorage.Metric) (memstorage.Metric, error)
	GetMetric(string) (memstorage.Metric, error)
	GetAllMetrics() memstorage.Metrics
}
