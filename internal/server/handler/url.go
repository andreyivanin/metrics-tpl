package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"metrics-tpl/internal/server/storage"
)

func (h *Handler) MetricUpdate(w http.ResponseWriter, r *http.Request) {
	var metric storage.Metric

	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	mvalue := chi.URLParam(r, "mvalue")
	w.Header().Set("Content-Type", "text/html")

	switch mtype {
	case "gauge":
		mvalueconv, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			return
		}

		metric = storage.Gauge(mvalueconv)

	case "counter":
		mvalueconv, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad metric value"))
			return
		}
		metric = storage.Counter(mvalueconv)

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad metric type"))
		return
	}

	h.Storage.UpdateMetric(mname, mtype, metric)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The metric " + mname + " was updated"))

}

func (h *Handler) MetricGet(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "mtype")
	mname := chi.URLParam(r, "mname")
	w.Header().Set("Content-Type", "text/html")

	metric, err := h.Storage.GetMetric(mname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The metric isn't found"))
		return
	}
	w.WriteHeader(http.StatusOK)

	switch mtype {
	case "gauge":
		if metric, ok := metric.(storage.Gauge); ok {
			metricconv := fmt.Sprintf("%.9g", metric)
			w.Write([]byte(metricconv))
			return
		}

	case "counter":
		if metric, ok := metric.(storage.Counter); ok {
			metricconv := strconv.Itoa(int(metric))
			w.Write([]byte(metricconv))
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Bad metric type"))
}

func (h *Handler) MetricSummary(w http.ResponseWriter, r *http.Request) {
	metrics := h.Storage.GetAllMetrics()
	w.Header().Set("Content-Type", "text/html")
	for name, metric := range metrics {
		valuestring := fmt.Sprintf("%v", metric)
		w.Write([]byte(name + ": " + valuestring + "\n"))
	}
}
