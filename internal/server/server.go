package server

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"metrics-tpl/internal/server/handler"
	"metrics-tpl/internal/server/middleware"
	"metrics-tpl/internal/server/storage"
)

func NewRouter(storage *storage.MemStorage) (chi.Router, error) {
	customHandler := handler.NewHandler(storage)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	// r.Use(chiMiddleware.Logger)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.GzipHandle)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", customHandler.MetricJSON)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", customHandler.MetricUpdate)
			r.Get("/", customHandler.MetricUpdate)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", customHandler.MetricSummaryJSON)
		r.Route("/{mtype}/{mname}", func(r chi.Router) {
			r.Get("/", customHandler.MetricGet)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", customHandler.MetricSummary)
	})

	return r, nil
}
