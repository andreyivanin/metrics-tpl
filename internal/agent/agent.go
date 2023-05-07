package agent

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Gauge float64
type Counter int64

type Metric interface{}

type Monitor struct {
	values       runtime.MemStats
	pollCounter  int
	SrvAddr      string
	UpdateTicker *time.Ticker
	SendTicker   *time.Ticker
	Metrics      map[string]Metric
}

func NewMonitor(cfg Config) Monitor {
	return Monitor{
		SrvAddr:      cfg.Address,
		UpdateTicker: time.NewTicker(cfg.PollInterval),
		SendTicker:   time.NewTicker(cfg.ReportInterval),
		Metrics:      make(map[string]Metric, 29),
	}
}

func (m *Monitor) UpdateMetrics() {
	runtime.ReadMemStats(&m.values)
	m.pollCounter++

	m.Metrics["Alloc"] = Gauge(m.values.Alloc)
	m.Metrics["BuckHashSys"] = Gauge(m.values.BuckHashSys)
	m.Metrics["Frees"] = Gauge(m.values.Frees)
	m.Metrics["GCCPUFraction"] = Gauge(m.values.GCCPUFraction)
	m.Metrics["GCSys"] = Gauge(m.values.GCSys)
	m.Metrics["HeapAlloc"] = Gauge(m.values.HeapAlloc)
	m.Metrics["HeapIdle"] = Gauge(m.values.HeapIdle)
	m.Metrics["HeapInuse"] = Gauge(m.values.HeapInuse)
	m.Metrics["HeapObjects"] = Gauge(m.values.HeapObjects)
	m.Metrics["HeapReleased"] = Gauge(m.values.HeapReleased)
	m.Metrics["HeapSys"] = Gauge(m.values.HeapSys)
	m.Metrics["LastGC"] = Gauge(m.values.LastGC)
	m.Metrics["Lookups"] = Gauge(m.values.Lookups)
	m.Metrics["MCacheInuse"] = Gauge(m.values.MCacheInuse)
	m.Metrics["MCacheSys"] = Gauge(m.values.MCacheSys)
	m.Metrics["MSpanInuse"] = Gauge(m.values.MSpanInuse)
	m.Metrics["MSpanSys"] = Gauge(m.values.MSpanSys)
	m.Metrics["Mallocs"] = Gauge(m.values.Mallocs)
	m.Metrics["NextGC"] = Gauge(m.values.NextGC)
	m.Metrics["NumForcedGC"] = Gauge(m.values.NumForcedGC)
	m.Metrics["NumGC"] = Gauge(m.values.NumGC)
	m.Metrics["OtherSys"] = Gauge(m.values.OtherSys)
	m.Metrics["PauseTotalNs"] = Gauge(m.values.PauseTotalNs)
	m.Metrics["StackInuse"] = Gauge(m.values.StackInuse)
	m.Metrics["StackSys"] = Gauge(m.values.StackSys)
	m.Metrics["Sys"] = Gauge(m.values.Sys)
	m.Metrics["TotalAlloc"] = Gauge(m.values.TotalAlloc)
	m.Metrics["RandomValue"] = Gauge(rand.Intn(100))
	m.Metrics["PollCount"] = Counter(m.pollCounter)
}

func (m *Monitor) SendMetrics() {
	client := http.Client{}

	for name, value := range m.Metrics {
		url := CreateURL(name, value, m.SrvAddr)
		log.Println(url)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		request.Header.Set("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
		}

		if response != nil {
			fmt.Println("Status code", response.Status)

			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(body))
		}
	}
}

func CreateURL(name string, value Metric, srvaddr string) string {
	var u url.URL
	var mtype string
	var valuestring string
	u.Scheme = "http"
	u.Host = srvaddr

	switch value := value.(type) {
	case Gauge:
		valuestring = fmt.Sprintf("%.f", value)
		mtype = "gauge"
	case Counter:
		valuestring = strconv.Itoa(int(value))
		mtype = "counter"
	}

	url := u.JoinPath("update", mtype, name, valuestring)
	return url.String()
}
