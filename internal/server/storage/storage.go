package storage

import (
	"context"
	"errors"
	"metrics-tpl/internal/server/config"
)

type Gauge float64
type Counter int64

type Metric interface{}
type Metrics map[string]Metric

type MemStorage struct {
	Metrics  map[string]Metric
	config   config.Config
	syncMode bool
}

func New(cfg config.Config) (*MemStorage, error) {
	return &MemStorage{
		Metrics: make(Metrics),
		config:  cfg,
	}, nil
}

func (s *MemStorage) ApplyConfig() error {
	if s.config.StoreFile != " " {
		s.enableFileStore()
	}

	if s.config.RestoreSavedData {
		s.Restore()
	}

	return nil
}

func (s *MemStorage) enableFileStore() {
	if s.config.StoreInterval == 0 {
		s.syncMode = true
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.SaveTicker(ctx, s.config.StoreInterval)
}

func (s *MemStorage) UpdateMetric(name string, m Metric) (Metric, error) {
	switch m.(type) {
	case Gauge:
		s.Metrics[name] = m
	case Counter:
		if existingMetric, ok := s.Metrics[name]; ok {
			updated := existingMetric.(Counter) + m.(Counter)
			s.Metrics[name] = updated
		} else {
			s.Metrics[name] = m
		}
	default:
		return nil, errors.New("the metric isn't found")
	}

	if s.syncMode {
		s.Save()
	}

	return s.Metrics[name], nil
}

func (s *MemStorage) GetMetric(mname string) (Metric, error) {
	if metric, ok := s.Metrics[mname]; ok {
		return metric, nil
	}

	return nil, errors.New("the metric isn't found")
}

func (s *MemStorage) GetAllMetrics() map[string]Metric {
	return s.Metrics

}
