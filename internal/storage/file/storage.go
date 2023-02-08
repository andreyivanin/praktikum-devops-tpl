package file

import (
	"bufio"
	"devops-tpl/internal/storage/memory"
	"encoding/json"
	"os"
)

type FileStorage struct {
	memory.MemStorage
}

var DB = New()

func New() *FileStorage {
	var d FileStorage
	d.GMetrics = make(map[string]memory.GaugeMetric)
	d.CMetrics = make(map[string]*memory.CounterMetric)
	return &d
}

func (s *FileStorage) UpdateGMetric(g memory.GaugeMetric) {
	s.GMetrics[g.Name] = g
}

func (s *FileStorage) UpdateCMetric(c memory.CounterMetric) {
	if existingMetric, ok := s.CMetrics[c.Name]; ok {
		existingMetric.Value = existingMetric.Value + c.Value
	} else {
		s.CMetrics[c.Name] = &c
	}
	// MetricUpdated <- true
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
	err := w.writer.Encode(memory.DB)
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
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &fileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *fileReader) ReadDatabase() (*FileStorage, error) {
	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	data := r.scanner.Bytes()

	DB := FileStorage{}
	err := json.Unmarshal(data, &DB)
	if err != nil {
		return nil, err
	}

	return &DB, nil
}
