package api

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	// 调用父结构体的方法
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		start := time.Now()
		handler.ServeHTTP(rec, r)
		duration := time.Since(start)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", r.Method).
			Dur("duration", duration).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Msg("received a HTTP request")
	})
}
