package storage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

type MetricFile struct {
	ID    string   `json:"id"`
	MType string   `json:"mtype"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}

func (s *MemStorage) Save() error {

	writer, err := NewWriter(s.config.StoreFile)
	if err != nil {
		log.Fatal(err)
	}

	defer writer.Close()

	MetricsFile := []MetricFile{}

	for name, metric := range s.Metrics {
		switch metric := metric.(type) {
		case Gauge:
			MetricsFile = append(MetricsFile, MetricFile{
				ID:    name,
				MType: "gauge",
				Value: (*float64)(&metric),
			})
		case Counter:
			MetricsFile = append(MetricsFile, MetricFile{
				ID:    name,
				MType: "counter",
				Delta: (*int64)(&metric),
			})
		}
	}

	err = writer.encoder.Encode(MetricsFile)
	if err != nil {
		return err
	}
	return nil

}

func (s *MemStorage) Restore() {
	reader, err := NewReader(s.config.StoreFile)
	if err != nil {
		log.Fatal(err)
	}

	checkFile, err := os.Stat(s.config.StoreFile)
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
		s.Metrics = restoredMetrics
	}
}

func (s *MemStorage) SaveTicker(ctx context.Context, storeint time.Duration) {
	ticker := time.NewTicker(storeint)
	defer ticker.Stop()

	for range ticker.C {
		s.Save()
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

func (r *fileReader) ReadDatabase() (Metrics, error) {
	MetricsFile := []MetricFile{}

	if err := r.reader.Decode(&MetricsFile); err != nil {
		return nil, err
	}

	Metrics := Metrics{}

	for _, metric := range MetricsFile {
		switch metric.MType {
		case "gauge":
			Metrics[metric.ID] = Gauge(*metric.Value)
		case "counter":
			Metrics[metric.ID] = Counter(*metric.Delta)
		}
	}
	return Metrics, nil
}
