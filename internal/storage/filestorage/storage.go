package filestorage

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"devops-tpl/internal/storage/memstorage"
)

type MetricFile struct {
	ID    string   `json:"id"`
	MType string   `json:"mtype"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}

type FileStorage struct {
	*memstorage.MemStorage
	storefile string
	SyncMode  bool
}

// var DB = New()

func New(storefile string) *FileStorage {
	memstorage := memstorage.New()
	return &FileStorage{
		MemStorage: memstorage,
		storefile:  storefile,
	}
}

// func (s *FileStorage) UpdateMetric(m memstorage.Metric) {
// 	s.Mu.Lock()
// 	defer s.Mu.Unlock()
// 	s.MemStorage.Metrics[m.Name] = m
// 	s.MemStorage.UpdateMetric(m)
// 	s.Save()
// }

// func (s *FileStorage) UpdateCMetric(c memstorage.CounterMetric) {
// 	s.Mu.Lock()
// 	defer s.Mu.Unlock()
// 	if existingMetric, ok := s.MemStorage.CMetrics[c.Name]; ok {
// 		existingMetric.Value = existingMetric.Value + c.Value
// 	} else {
// 		s.MemStorage.CMetrics[c.Name] = &c
// 	}
// 	s.Save()
// }

func (s *FileStorage) UpdateMetric(name string, m memstorage.Metric) (memstorage.Metric, error) {
	s.MemStorage.UpdateMetric(name, m)
	s.Save()
	return s.Metrics[name], nil
}

func (s *FileStorage) Save() error {

	writer, err := NewWriter(s.storefile)
	if err != nil {
		log.Fatal(err)
	}

	defer writer.Close()

	MetricsFile := []MetricFile{}

	for name, metric := range s.MemStorage.Metrics {
		switch metric := metric.(type) {
		case memstorage.GaugeMetric:
			MetricsFile = append(MetricsFile, MetricFile{
				ID:    name,
				MType: "gauge",
				Value: &metric.Value,
			})
		case memstorage.CounterMetric:
			MetricsFile = append(MetricsFile, MetricFile{
				ID:    name,
				MType: "counter",
				Delta: &metric.Delta,
			})
		}
	}

	err = writer.encoder.Encode(MetricsFile)
	if err != nil {
		return err
	}
	return nil

}

func (s *FileStorage) Restore(storefile string) {
	reader, err := NewReader(storefile)
	if err != nil {
		log.Fatal(err)
	}

	checkFile, err := os.Stat(storefile)
	if err != nil {
		log.Fatal(err)
	}

	size := checkFile.Size()

	if size == 0 {
		s.Save()
	}

	if restoredMetrics, err := reader.ReadDatabase(); err != nil {
		log.Fatal(err)
	} else {
		s.MemStorage.Metrics = *restoredMetrics
	}
}

func (s *FileStorage) SaveTicker(storeint time.Duration) {
	ticker := time.NewTicker(storeint)
	defer ticker.Stop()

	for range ticker.C {
		s.Save()
	}
}

func (s *FileStorage) GetAllMetrics() memstorage.Metrics {
	return s.MemStorage.GetAllMetrics()
}

// func (s *FileStorage) SaveTicker(ctx context.Context, storeint time.Duration) {
// 	ticker := time.NewTicker(storeint)
// 	for range ticker.C {
// 		s.Save()
// 	}

// }

type fileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(filename string) (*fileWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *fileWriter) Close() error {
	return w.file.Close()
}

type fileReader struct {
	file   *os.File
	reader *json.Decoder
}

func NewReader(filename string) (*fileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &fileReader{
		file:   file,
		reader: json.NewDecoder(file),
	}, nil
}

func (r *fileReader) ReadDatabase() (*memstorage.Metrics, error) {
	MetricsFile := []MetricFile{}

	if err := r.reader.Decode(&MetricsFile); err != nil {
		return nil, err
	}

	Metrics := memstorage.Metrics{}

	for _, metric := range MetricsFile {
		switch metric.MType {
		case "gauge":
			Metrics[metric.ID] = memstorage.GaugeMetric{
				MType: metric.MType,
				Value: *metric.Value,
			}
		case "counter":
			Metrics[metric.ID] = memstorage.CounterMetric{
				MType: metric.MType,
				Delta: *metric.Delta,
			}
		}
	}

	return &Metrics, nil
}
