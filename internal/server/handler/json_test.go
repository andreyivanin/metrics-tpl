package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"metrics-tpl/internal/server/config"
	"metrics-tpl/internal/server/storage"
)

func TestHandler_MetricUpdateJSON(t *testing.T) {

	type request struct {
		url  string
		body string
	}

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name string
		req  request
		want want
	}{
		{
			name: "test#1: good gauge",
			req: request{
				url:  "/update",
				body: `{"id":"HeapReleased","type":"gauge","value":2940930}`,
			},
			want: want{
				code: 200,
				body: `{"id":"HeapReleased","type":"gauge","value":2940930}`,
			},
		},
		{
			name: "test#2: good counter",
			req: request{
				url:  "/update",
				body: `{"id":"RandomValue","type":"counter","delta":47}`,
			},
			want: want{
				code: 200,
				body: `{"id":"RandomValue","type":"counter","delta":47}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, _ := storage.New(config.Config{})

			rBody := bytes.NewReader([]byte(tt.req.body))

			r := httptest.NewRequest(http.MethodGet, tt.req.url, rBody)

			w := httptest.NewRecorder()

			handler := NewHandler(storage)

			handler.MetricUpdateJSON(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)

			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.JSONEq(t, tt.want.body, string(resultBody))

		})
	}
}
