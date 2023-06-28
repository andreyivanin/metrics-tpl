package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type signWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (w signWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Sign(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedSign := r.Header.Get("HashSHA256")
			if receivedSign == "" {
				next.ServeHTTP(w, r)
				return
			}

			b, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body.Close()

			r.Body = io.NopCloser(bytes.NewBuffer(b))

			reqSign := CreateSign(b, []byte(key))

			if receivedSign != hex.EncodeToString(reqSign) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			sW := signWriter{}

			respSign := CreateSign((sW.body.Bytes()), []byte(key))

			w.Header().Set("HashSHA256", hex.EncodeToString(respSign))

			next.ServeHTTP(signWriter{ResponseWriter: w}, r)
		})
	}
}

func CreateSign(payload, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(payload)
	return h.Sum(nil)
}
