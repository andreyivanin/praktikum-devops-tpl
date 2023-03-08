package storage

import (
	"devops-tpl/internal/storage/memstorage"
)

type Storage interface {
	UpdateGMetric(memstorage.GaugeMetric)
	UpdateCMetric(memstorage.CounterMetric)
	GetGMetric(mname string) (memstorage.GaugeMetric, error)
	GetCMetric(mname string) (*memstorage.CounterMetric, error)
	GetStorage() *memstorage.MemStorage
}
