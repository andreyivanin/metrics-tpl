package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateMetric(t *testing.T) {

	type fields struct {
		name   string
		mtype  string
		metric Metric
	}

	tests := []struct {
		name   string
		metric fields
		want   MemStorage
	}{
		{
			name: "update gauge metric",
			metric: fields{
				name:   "Alloc",
				mtype:  _GAUGE,
				metric: Gauge(1223113),
			},
			want: MemStorage{
				Metrics: Metrics{"Alloc": Gauge(1223113)},
			},
		},
		{
			name: "update counter metric",
			metric: fields{
				name:   "RandomValue",
				mtype:  _COUNTER,
				metric: Counter(67),
			},
			want: MemStorage{
				Metrics: Metrics{"RandomValue": Counter(134)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var DB = MemStorage{
				Metrics: make(Metrics),
			}
			DB.UpdateMetric(tt.metric.name, tt.metric.mtype, tt.metric.metric)
			DB.UpdateMetric(tt.metric.name, tt.metric.mtype, tt.metric.metric)
			assert.Equal(t, tt.want, DB)
		})
	}
}
