package storage

import (
	"context"
	"errors"
	"log"
	"sync"

	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/models"
)

const (
	_GAUGE   = "gauge"
	_COUNTER = "counter"
)

type MemStorage struct {
	Metrics  models.Metrics
	Mu       *sync.Mutex
	config   config.Config
	syncMode bool
}

func newMemStorage(ctx context.Context, cfg config.Config) (*MemStorage, error) {
	memStorage := &MemStorage{
		Metrics: make(models.Metrics),
		Mu:      new(sync.Mutex),
		config:  cfg,
	}

	if memStorage.config.StoreFile != " " {
		memStorage.enableFileStore(ctx)
	}

	if memStorage.config.RestoreSavedData {
		err := memStorage.Restore()
		if err != nil {
			log.Print(err)
		}
	}

	return memStorage, nil

}

func (s *MemStorage) enableFileStore(ctx context.Context) {
	if s.config.StoreInterval == 0 {
		s.syncMode = true
		return
	}

	go s.SaveTicker(ctx, s.config.StoreInterval)
}

func (s *MemStorage) enableSQLStore(ctx context.Context) {
	if s.config.DatabaseDSN == "" {
		return
	}

	go s.SaveTicker(ctx, s.config.StoreInterval)
}

func (s *MemStorage) UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	switch mtype {
	case _GAUGE:
		s.Metrics[name] = m
	case _COUNTER:
		if existingMetric, ok := s.Metrics[name]; ok {
			updated := existingMetric.(models.Counter) + m.(models.Counter)
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

func (s *MemStorage) GetMetric(ctx context.Context, mname string) (models.Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if metric, ok := s.Metrics[mname]; ok {
		return metric, nil
	}

	return nil, ErrNotFound
}

func (s *MemStorage) GetAllMetrics(ctx context.Context) (models.Metrics, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return s.Metrics, nil
}

func (s *MemStorage) GetConfig() config.Config {
	return s.config
}
