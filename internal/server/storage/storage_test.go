package storage

import (
	"context"
	"sync"
	"testing"

	"metrics-tpl/internal/server/models"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateMetric(t *testing.T) {

	type fields struct {
		name   string
		mtype  string
		metric models.Metric
	}

	tests := []struct {
		name   string
		metric fields
		want   models.Metrics
	}{
		{
			name: "update gauge metric",
			metric: fields{
				name:   "Alloc",
				mtype:  _GAUGE,
				metric: models.Gauge(1223113),
			},
			want: models.Metrics{"Alloc": models.Gauge(1223113)},
		},
		{
			name: "update counter metric",
			metric: fields{
				name:   "RandomValue",
				mtype:  _COUNTER,
				metric: models.Counter(67),
			},
			want: models.Metrics{"RandomValue": models.Counter(134)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var DB = MemStorage{
				Metrics: make(models.Metrics),
				Mu:      new(sync.Mutex),
			}
			DB.UpdateMetric(ctx, tt.metric.name, tt.metric.mtype, tt.metric.metric)
			DB.UpdateMetric(ctx, tt.metric.name, tt.metric.mtype, tt.metric.metric)
			assert.Equal(t, tt.want, DB.Metrics)
		})
	}
}
