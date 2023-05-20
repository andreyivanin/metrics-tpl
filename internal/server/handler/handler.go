package handler

import (
	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/storage"
)

type Repository interface {
	UpdateMetric(name, mtype string, m storage.Metric) (storage.Metric, error)
	GetMetric(mname string) (storage.Metric, error)
	GetAllMetrics() storage.Metrics
	GetConfig() config.Config
}

type Handler struct {
	Storage Repository
}

func NewHandler(storage Repository) *Handler {
	return &Handler{storage}
}
