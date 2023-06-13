package storage

import (
	"context"
	"errors"

	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/models"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error)
	GetMetric(ctx context.Context, mname string) (models.Metric, error)
	GetAllMetrics(ctx context.Context) (models.Metrics, error)
	GetConfig() config.Config
}

func New(ctx context.Context, cfg config.Config) (Storage, error) {
	if cfg.DatabaseDSN == "" {
		return newMemStorage(ctx, cfg)
	}

	return newSQL(ctx, cfg)
}
