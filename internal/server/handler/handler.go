package handler

import (
	"metrics-tpl/internal/server/storage"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Storage *storage.MemStorage
}

func NewHandler(storage *storage.MemStorage) *Handler {
	return &Handler{storage}
}

func (h *Handler) MetricUpdate(w http.ResponseWriter, r *http.Request) {
	var metric storage.Metric

	url := r.URL.Path
	fields := strings.Split(url, "/")

	// mtype := chi.URLParam(r, "mtype")
	// mname := chi.URLParam(r, "mname")
	// mvalue := chi.URLParam(r, "mvalue")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// switch mtype {
	// case "gauge":
	// 	mvalueconv, err := strconv.ParseFloat(mvalue, 64)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte("Bad metric value"))
	// 		return
	// 	}

	// 	metric = storage.Gauge(mvalueconv)

	// case "counter":
	// 	mvalueconv, err := strconv.ParseInt(mvalue, 10, 64)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte("Bad metric value"))
	// 		return
	// 	}
	// 	metric = storage.Counter(mvalueconv)

	// default:
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Bad metric type"))
	// 	return
	// }

	if len(fields) == 5 && fields[1] == "update" {
		switch fields[2] {
		case "gauge":
			mvalueconv, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad metric value"))
				return
			}

			metric = storage.Gauge(mvalueconv)

		case "counter":
			mvalueconv, err := strconv.ParseInt(fields[4], 10, 64)
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

		h.Storage.UpdateMetric(fields[3], metric)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The metric " + fields[3] + " was updated"))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad metric type"))
		return
	}

	// h.Storage.UpdateMetric(mname, metric)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("The metric " + mname + " was updated"))

}
