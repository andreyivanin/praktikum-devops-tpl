package filestorage

import (
	"bufio"
	"devops-tpl/internal/storage/memstorage"
	"encoding/json"
	"log"
	"os"
	"time"
)

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

func (s *FileStorage) UpdateGMetric(g memstorage.GaugeMetric) {
	s.MemStorage.GMetrics[g.Name] = g
	s.Save()
}

func (s *FileStorage) UpdateCMetric(c memstorage.CounterMetric) {
	if existingMetric, ok := s.MemStorage.CMetrics[c.Name]; ok {
		existingMetric.Value = existingMetric.Value + c.Value
	} else {
		s.MemStorage.CMetrics[c.Name] = &c
	}
	s.Save()
}

func (s *FileStorage) Save() error {
	writer, err := NewWriter(s.storefile)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.encoder.Encode(s.MemStorage)
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

	if restored, err := reader.ReadDatabase(); err != nil {
		log.Fatal(err)
	} else {
		s.MemStorage = restored.MemStorage
	}
}

func (s *FileStorage) SaveTicker(storeint time.Duration) {
	ticker := time.NewTicker(storeint)
	for range ticker.C {
		s.Save()
	}
}

// func (s *FileStorage) SaveTicker(ctx context.Context, storeint time.Duration) {
// 	ticker := time.NewTicker(storeint)
// 	for range ticker.C {
// 		s.Save()
// 	}

// }

func (s *FileStorage) GetStorage() *memstorage.MemStorage {
	return &memstorage.MemStorage{
		GMetrics: s.GMetrics,
		CMetrics: s.CMetrics,
	}
}

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
