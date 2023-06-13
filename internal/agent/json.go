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
	"syscall"
	"time"
)

const (
	_GAUGE   = "gauge"
	_COUNTER = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Monitor) SendMetricsJSON() error {
	url := CreateURLJSON(m.SrvAddr, "update")
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := func(ctx context.Context) error {
		for name, metric := range m.Metrics {
			jsonMetric := Metrics{}

			switch metric := metric.(type) {
			case Gauge:
				jsonMetric = Metrics{
					ID:    name,
					MType: _GAUGE,
					Value: (*float64)(&metric),
				}
			case Counter:
				jsonMetric = Metrics{
					ID:    name,
					MType: _COUNTER,
					Delta: (*int64)(&metric),
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

func (m *Monitor) SendMetricsGroupJSON() error {
	url := CreateURLJSON(m.SrvAddr, "updates")
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	jsonMetrics := make([]Metrics, 0, 29)

	err := func(ctx context.Context) error {
		for name, metric := range m.Metrics {
			jsonMetric := Metrics{}

			switch metric := metric.(type) {
			case Gauge:
				jsonMetric = Metrics{
					ID:    name,
					MType: _GAUGE,
					Value: (*float64)(&metric),
				}
			case Counter:
				jsonMetric = Metrics{
					ID:    name,
					MType: _COUNTER,
					Delta: (*int64)(&metric),
				}
			}

			jsonMetrics = append(jsonMetrics, jsonMetric)
		}

		return nil
	}(ctx)
	if err != nil {
		return err
	}

	bytesMetrics, err := json.Marshal(jsonMetrics)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(bytesMetrics)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := sendRequest(&client, request)
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
	if err != nil {
		return err
	}

	return nil
}

func sendRequest(client *http.Client, request *http.Request) (response *http.Response, err error) {
	var attemptsCount int
	seconds := []int{1, 3, 5}

	waitDuration := func(seconds []int) []time.Duration {
		duration := make([]time.Duration, len(seconds))
		for i, sec := range seconds {
			duration[i] = time.Duration(sec) * time.Second
		}
		return duration
	}(seconds)

	for attemptsCount < _MAXSENDATTEMPTS {
		response, err = client.Do(request)
		if errors.Is(err, syscall.ECONNREFUSED) {
			time.Sleep(waitDuration[attemptsCount])
			attemptsCount++
			continue
		}

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, err
}

func CreateURLJSON(srv, path string) string {
	var u url.URL
	u.Scheme = _PROTOCOL
	u.Host = srv
	url := u.JoinPath(path)
	return url.String()
}
