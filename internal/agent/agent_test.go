package agent

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {

	tests := []struct {
		name    string
		metrics map[string]Metric
		want    string
	}{
		{
			name:    "good test: Gaugemetric",
			metrics: map[string]Metric{"Alloc": Gauge(150)},
			want:    "/update/gauge/Alloc/150",
		},
		{
			name:    "good test: Countermetric",
			metrics: map[string]Metric{"PollCount": Counter(55)},
			want:    "/update/counter/PollCount/55",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mon := Monitor{
				Metrics: tt.metrics,
			}

			l, err := net.Listen("tcp", "127.0.0.1:8080")
			if err != nil {
				log.Fatal(err)
			}

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				path := req.URL.Path
				assert.Equal(t, tt.want, path)
			}))
			defer func() { ts.Close() }()

			ts.Listener.Close()
			ts.Listener = l
			ts.Start()

			mon.SendMetrics()

		})
	}
}
