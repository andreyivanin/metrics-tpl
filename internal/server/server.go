package server

import (
	"context"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/handler"
	"metrics-tpl/internal/server/middleware"
	"metrics-tpl/internal/server/models"
)

type Storage interface {
	UpdateMetric(ctx context.Context, name, mtype string, m models.Metric) (models.Metric, error)
	GetMetric(ctx context.Context, mname string) (models.Metric, error)
	GetAllMetrics(ctx context.Context) (models.Metrics, error)
	GetConfig() config.Config
}

func NewRouter(storage Storage, cfg config.Config) chi.Router {
	customHandler := handler.NewHandler(storage)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	// r.Use(chiMiddleware.Logger)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.GzipHandle)
	if cfg.Key != "" {
		r.Use(middleware.Sign(cfg.Key))
	}

	r.Route("/update", func(r chi.Router) {
		r.Post("/", customHandler.MetricUpdateJSON)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", customHandler.MetricUpdate)
			r.Get("/", customHandler.MetricUpdate)
		})
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", customHandler.MetricsGroupUpdateJSON)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", customHandler.MetricGetJSON)
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", customHandler.MetricGet)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", customHandler.MetricSummary)
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", customHandler.TestDBConnection)
	})

	return r
}
