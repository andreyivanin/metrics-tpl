package storage

import "errors"

type Gauge float64
type Counter int64

type Metric interface{}

type MemStorage struct {
	Metrics map[string]Metric
}

func New() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]Metric),
	}
}

func (s *MemStorage) UpdateMetric(name string, m Metric) (Metric, error) {
	switch m.(type) {
	case Gauge:
		s.Metrics[name] = m
	case Counter:
		if existingMetric, ok := s.Metrics[name]; ok {
			updated := existingMetric.(Counter) + m.(Counter)
			s.Metrics[name] = updated
		} else {
			s.Metrics[name] = m
		}
	default:
		return nil, errors.New("the metric isn't found")
	}

	return s.Metrics[name], nil
}

func (s *MemStorage) GetMetric(mname string) (Metric, error) {
	if metric, ok := s.Metrics[mname]; ok {
		return metric, nil
	}

	return nil, errors.New("the metric isn't found")
}

func (s *MemStorage) GetAllMetrics() map[string]Metric {
	return s.Metrics

}
