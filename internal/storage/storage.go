package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

var MetricUpdated = make(chan bool)

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

var DB = createDB()

func createDB() *MemStorage {
	var d MemStorage
	d.GMetrics = make(map[string]GaugeMetric)
	d.CMetrics = make(map[string]*CounterMetric)
	return &d
}

func UpdateGMetric(g GaugeMetric, s *MemStorage) {
	s.GMetrics[g.Name] = g
	MetricUpdated <- true

}

func UpdateCMetric(c CounterMetric, s *MemStorage) {
	if existingMetric, ok := s.CMetrics[c.Name]; ok {
		existingMetric.Value = existingMetric.Value + c.Value
	} else {
		s.CMetrics[c.Name] = &c
	}
	MetricUpdated <- true
}

func GetGMetric(mname string) (GaugeMetric, error) {
	if metric, ok := DB.GMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func GetCMetric(mname string) (*CounterMetric, error) {
	if metric, ok := DB.CMetrics[mname]; ok {
		return metric, nil
	} else {
		err := errors.New("the metric isn't found")
		return metric, err
	}
}

func GetMetricSummary() *MemStorage {
	return DB
}

type fileWriter struct {
	file   *os.File
	writer *json.Encoder
}

func NewWriter(filename string) (*fileWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file:   file,
		writer: json.NewEncoder(file),
	}, nil
}

func (w *fileWriter) WriteDatabase() error {
	err := w.writer.Encode(DB)
	if err != nil {
		return err
	}
	return nil
}

func (w *fileWriter) Close() error {
	return w.file.Close()
}

// type fileReader struct {
// 	file   *os.File
// 	reader *json.Decoder
// }

// func NewReader(filename string) (*fileReader, error) {
// 	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &fileReader{
// 		file:   file,
// 		reader: json.NewDecoder(file),
// 	}, nil
// }

// func (r *fileReader) ReadDatabase() error {
// 	if err := r.reader.Decode(&DB); err != nil {
// 		return err
// 	}
// 	return nil
// }

type fileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewReader(filename string) (*fileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	return &fileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *fileReader) ReadDatabase() (*MemStorage, error) {
	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	data := r.scanner.Bytes()

	DB := MemStorage{}
	err := json.Unmarshal(data, &DB)
	if err != nil {
		return nil, err
	}

	return &DB, nil
}
