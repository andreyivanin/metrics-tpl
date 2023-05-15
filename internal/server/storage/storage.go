package storage

import (
	"context"
	"errors"
	"log"
	"sync"

	"metrics-tpl/internal/server/config"
)

const (
	_GAUGE   = "gauge"
	_COUNTER = "counter"
)

var ErrNotFound = errors.New("not found")

type Gauge float64
type Counter int64

type Metric interface{}
type Metrics map[string]Metric

type MemStorage struct {
	Metrics  map[string]Metric
	Mu       *sync.Mutex
	config   config.Config
	syncMode bool
}

func New(cfg config.Config) (*MemStorage, error) {
	memStorage := &MemStorage{
		Metrics: make(Metrics),
		Mu:      new(sync.Mutex),
		config:  cfg,
	}

	if memStorage.config.StoreFile != " " {
		memStorage.enableFileStore()
	}

	if memStorage.config.RestoreSavedData {
		err := memStorage.Restore()
		if err != nil {
			log.Print(err)
		}
	}

	return memStorage, nil

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

func (s *MemStorage) UpdateMetric(name, mtype string, m Metric) (Metric, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	switch mtype {
	case _GAUGE:
		s.Metrics[name] = m
	case _COUNTER:
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

	return nil, ErrNotFound
}

func (s *MemStorage) GetAllMetrics() map[string]Metric {
	return s.Metrics

}
