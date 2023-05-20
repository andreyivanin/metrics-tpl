package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"metrics-tpl/internal/server/storage"
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

func (h *Handler) MetricUpdateJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metric storage.Metric

	switch jsonmetric.MType {
	case _GAUGE:
		metric = storage.Gauge(*jsonmetric.Value)
	case _COUNTER:
		metric = storage.Counter(*jsonmetric.Delta)
	}

	updatedMetric, err := h.Storage.UpdateMetric(jsonmetric.ID, jsonmetric.MType, metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch updatedMetric := updatedMetric.(type) {
	case storage.Gauge:
		jsonmetric.Value = (*float64)(&updatedMetric)

	case storage.Counter:
		jsonmetric.Delta = (*int64)(&updatedMetric)
	}

	metricsJSON, err := json.Marshal(jsonmetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}

func (h *Handler) MetricGetJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonmetric := Metrics{}

	if err := json.Unmarshal(b, &jsonmetric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Storage.GetMetric(jsonmetric.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("One or several metrics weren't found"))
		return
	}

	switch jsonmetric.MType {
	case _GAUGE:
		metric := metric.(storage.Gauge)
		jsonmetric.Value = (*float64)(&metric)

	case _COUNTER:
		metric := metric.(storage.Counter)
		jsonmetric.Delta = (*int64)(&metric)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("wrong metric type"))
		log.Println("wrong metric type")
	}

	metricsJSON, err := json.Marshal(jsonmetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(metricsJSON)
}
