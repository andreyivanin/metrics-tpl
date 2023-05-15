package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	_GAUGE   = "gauge"
	_COUNTER = "counter"
)

type Metrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Monitor) SendMetricsJSON() error {
	url := CreateURLJSON(m.SrvAddr)
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := func(ctx context.Context) error {
		for name, metric := range m.Metrics {
			jsonMetric := Metrics{}

			switch metric := metric.(type) {
			case Gauge:
				jsonMetric = Metrics{
					ID:    name,
					MType: _GAUGE,
					Value: (float64)(metric),
				}
			case Counter:
				jsonMetric = Metrics{
					ID:    name,
					MType: _COUNTER,
					Delta: (int64)(metric),
				}

			}
			bytesMetric, err := json.Marshal(jsonMetric)
			if err != nil {
				return err
			}

			body := bytes.NewBuffer(bytesMetric)

			request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
			if err != nil {
				return err
			}

			request.Header.Set("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				return err
			}

			if response == nil {
				return errors.New("empty response")
			}

			log.Println("Status code", response.Status)

			defer response.Body.Close()

			bodyResp, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			log.Printf("Response body:\n %v\n", string(bodyResp))
		}
		return nil
	}(ctx)

	if err != nil {
		return err
	}

	return nil
}

func CreateURLJSON(srv string) string {
	var u url.URL
	u.Scheme = _PROTOCOL
	u.Host = srv
	url := u.JoinPath("update")
	return url.String()
}
