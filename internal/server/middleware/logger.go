package middleware

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar = *logger.Sugar()

	return nil
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			statusCode: 0,
			size:       0,
		}

		lw := logResponseWriter{
			ResponseWriter: w,
			resposeData:    responseData,
		}

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		sugar.Infoln(
			"URI:", uri,
			"Method:", method,
			"StatusCode:", responseData.statusCode,
			"Duration:", duration,
			"Size:", responseData.size,
		)
	})
}

type responseData struct {
	statusCode int
	size       int
}

type logResponseWriter struct {
	http.ResponseWriter
	resposeData *responseData
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.resposeData.size += size
	return size, err
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.resposeData.statusCode = statusCode
}

func init() {
	err := InitLogger()
	if err != nil {
		log.Println(err)
	}
}
