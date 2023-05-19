package handler

import (
	"context"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/middleware"
	"metrics-tpl/internal/server/storage"
)

func NewRouter(storage *storage.MemStorage) (chi.Router, error) {
	customHandler := NewHandler(storage)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	// r.Use(chiMiddleware.Logger)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.GzipHandle)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", customHandler.MetricUpdateJSON)
		r.Route("/{mtype}/{mname}/{mvalue}", func(r chi.Router) {
			r.Post("/", customHandler.MetricUpdate)
			r.Get("/", customHandler.MetricUpdate)
		})
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

	return r, nil
}

func testRequestURL(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestHandler_MetricUpdate(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "test#1: good gauge",
			url:  "/update/gauge/Metric/100",
			want: want{
				code: 200,
				body: `The metric Metric was updated`,
			},
		},
		{
			name: "test#2: good counter2",
			url:  "/update/counter/Metric/100",
			want: want{
				code: 200,
				body: `The metric Metric was updated`,
			},
		},
		{
			name: "test#3: bad metric type",
			url:  "/update/gaugecounter/Metric/100",
			want: want{
				code: 400,
				body: `Bad metric type`,
			},
		},
		{
			name: "test#4: bad url",
			url:  "/update/gaug100",
			want: want{
				code: 404,
				body: `404 page not found
`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			storage, _ := storage.New(ctx, config.Config{})
			router, _ := NewRouter(storage)
			ts := httptest.NewServer(router)
			defer ts.Close()

			code, body := testRequestURL(t, ts, http.MethodPost, tt.url)

			assert.Equal(t, tt.want.code, code)
			assert.Equal(t, tt.want.body, body)
		})
	}
}
