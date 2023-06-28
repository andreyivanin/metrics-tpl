package agent

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	ErrorRateExceed = errors.New("rate limit was exceeded")
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
	RateTicker   *time.Ticker
	Key          string
	Storage      *Storage
	RateLimiter  RateLimiter
}

type Storage struct {
	Metrics map[string]Metric
	Mu      *sync.Mutex
}

func newStorage() *Storage {
	return &Storage{
		Metrics: make(map[string]Metric, 32),
		Mu:      new(sync.Mutex),
	}
}

func NewMonitor(cfg Config) Monitor {
	storage := newStorage()
	rateLimiter := NewRateLimiter(2)

	return Monitor{
		SrvAddr:      cfg.Address,
		UpdateTicker: time.NewTicker(cfg.PollInterval),
		SendTicker:   time.NewTicker(cfg.ReportInterval),
		RateTicker:   time.NewTicker(10 * time.Second),
		Key:          cfg.Key,
		Storage:      storage,
		RateLimiter:  *rateLimiter,
	}
}

func (m *Monitor) UpdateMetrics() {
	m.Storage.Mu.Lock()
	defer m.Storage.Mu.Unlock()

	runtime.ReadMemStats(&m.values)
	m.pollCounter++

	m.Storage.Metrics["Alloc"] = Gauge(m.values.Alloc)
	m.Storage.Metrics["BuckHashSys"] = Gauge(m.values.BuckHashSys)
	m.Storage.Metrics["Frees"] = Gauge(m.values.Frees)
	m.Storage.Metrics["GCCPUFraction"] = Gauge(m.values.GCCPUFraction)
	m.Storage.Metrics["GCSys"] = Gauge(m.values.GCSys)
	m.Storage.Metrics["HeapAlloc"] = Gauge(m.values.HeapAlloc)
	m.Storage.Metrics["HeapIdle"] = Gauge(m.values.HeapIdle)
	m.Storage.Metrics["HeapInuse"] = Gauge(m.values.HeapInuse)
	m.Storage.Metrics["HeapObjects"] = Gauge(m.values.HeapObjects)
	m.Storage.Metrics["HeapReleased"] = Gauge(m.values.HeapReleased)
	m.Storage.Metrics["HeapSys"] = Gauge(m.values.HeapSys)
	m.Storage.Metrics["LastGC"] = Gauge(m.values.LastGC)
	m.Storage.Metrics["Lookups"] = Gauge(m.values.Lookups)
	m.Storage.Metrics["MCacheInuse"] = Gauge(m.values.MCacheInuse)
	m.Storage.Metrics["MCacheSys"] = Gauge(m.values.MCacheSys)
	m.Storage.Metrics["MSpanInuse"] = Gauge(m.values.MSpanInuse)
	m.Storage.Metrics["MSpanSys"] = Gauge(m.values.MSpanSys)
	m.Storage.Metrics["Mallocs"] = Gauge(m.values.Mallocs)
	m.Storage.Metrics["NextGC"] = Gauge(m.values.NextGC)
	m.Storage.Metrics["NumForcedGC"] = Gauge(m.values.NumForcedGC)
	m.Storage.Metrics["NumGC"] = Gauge(m.values.NumGC)
	m.Storage.Metrics["OtherSys"] = Gauge(m.values.OtherSys)
	m.Storage.Metrics["PauseTotalNs"] = Gauge(m.values.PauseTotalNs)
	m.Storage.Metrics["StackInuse"] = Gauge(m.values.StackInuse)
	m.Storage.Metrics["StackSys"] = Gauge(m.values.StackSys)
	m.Storage.Metrics["Sys"] = Gauge(m.values.Sys)
	m.Storage.Metrics["TotalAlloc"] = Gauge(m.values.TotalAlloc)
	m.Storage.Metrics["RandomValue"] = Gauge(rand.Intn(100))
	m.Storage.Metrics["PollCount"] = Counter(m.pollCounter)
}

func (m *Monitor) UpdateCustomMetrics() {
	m.Storage.Mu.Lock()
	defer m.Storage.Mu.Unlock()

	v, _ := mem.VirtualMemory()

	m.Storage.Metrics["TotalMemory"] = Gauge(v.Total)
	m.Storage.Metrics["FreeMemory"] = Gauge(v.Free)

	CPUUtil, err := cpu.Percent(0, true)
	if err != nil {
		return
	}

	m.Storage.Metrics["CPUutilization1"] = Gauge(CPUUtil[0])
}

func (m *Monitor) SendMetrics() error {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := func(ctx context.Context) error {
		for name, value := range m.Storage.Metrics {
			url := CreateURL(name, value, m.SrvAddr)
			log.Println(url)

			request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
			if err != nil {
				return err
			}

			request.Header.Set("Content-Type", "text/plain")
			response, err := client.Do(request)
			if err != nil {
				return err
			}

			if response == nil {
				return errors.New("empty response")
			}

			fmt.Println("Status code", response.Status)

			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(body))
		}

		return nil
	}(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m.NewRateLimiter()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {

		for {
			select {
			case <-ctx.Done():
				wg.Done()
			case <-m.UpdateTicker.C:
				m.UpdateMetrics()
				fmt.Println("Metrics update", " - ", time.Now())
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
			case <-m.UpdateTicker.C:
				m.UpdateCustomMetrics()
				fmt.Println("Custom Metrics update", " - ", time.Now())

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
			case <-m.SendTicker.C:
				if ok := <-m.RateLimiter.rateCh; !ok {
					log.Println(ErrorRateExceed)
					continue
				}

				m.SendMetricsGroupJSON(ctx)

				// if err != nil {
				// 	log.Print(err)
				// }
				m.RateLimiter.Release()
				fmt.Println("Metrics send", " - ", time.Now())
			}
		}
	}()

	wg.Wait()

	return nil
}

type RateLimiter struct {
	rateCh chan bool
}

func NewRateLimiter(maxReq int) *RateLimiter {
	return &RateLimiter{
		rateCh: make(chan bool, maxReq),
	}
}

func (rl *RateLimiter) Allocate() {
	rl.rateCh <- true
}

func (rl *RateLimiter) Release() {
	rl.rateCh <- false
}

func (m *Monitor) NewRateLimiter() {
	go func() {
		for {
			select {
			case <-m.RateTicker.C:
				m.RateLimiter.Allocate()
			}
		}
	}()
}

func CreateSign(payload, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(payload)
	return h.Sum(nil)
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
