package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Monitor) SendMetricsJSON() {
	url := CreateURLJSON(m.SrvAddr)
	client := http.Client{}

	for name, metric := range m.Metrics {
		jsonMetric := Metrics{}

		switch metric := metric.(type) {
		case Gauge:
			jsonMetric = Metrics{
				ID:    name,
				MType: "gauge",
				Value: (*float64)(&metric),
			}
		case Counter:
			jsonMetric = Metrics{
				ID:    name,
				MType: "counter",
				Delta: (*int64)(&metric),
			}

		}
		bytesMetric, err := json.Marshal(jsonMetric)
		if err != nil {
			panic(err)
		}

		body := bytes.NewBuffer(bytesMetric)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
		if err != nil {
			log.Print(err)
		}

		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			log.Print(err)
		}

		// requestDump, err := httputil.DumpRequest(request, true)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println(string(requestDump))

		if response != nil {
			fmt.Println("Status code", response.Status)

			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Response body:\n %v\n", string(body))
		}

	}
}

func CreateURLJSON(srv string) string {
	var u url.URL
	u.Scheme = PROTOCOL
	u.Host = srv
	url := u.JoinPath("update")
	return url.String()
}
