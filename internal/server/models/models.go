package models

type Gauge float64
type Counter int64

type Metric interface{}
type Metrics map[string]Metric
