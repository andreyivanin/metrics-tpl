package handler

import (
	"context"
	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/models"
	"time"
)

const timeoutCTX = 3 * time.Second

type Storage interface {
	UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error)
	GetMetric(ctx context.Context, mname string) (models.Metric, error)
	GetAllMetrics(ctx context.Context) (models.Metrics, error)
	GetConfig() config.Config
}

type Handler struct {
	Storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage}
}
