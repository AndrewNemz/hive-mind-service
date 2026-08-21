package middleware

import (
	"hiv_mind/pkg/logger"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lg := logger.Get().Sugar()
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		lg.Info(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	})
}

func ResponseInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg := logger.Get().Sugar()
		rd := &responseData{
			status: 0,
			size:   0,
		}
		lr := loggingResponseWriter{w, rd}

		next.ServeHTTP(&lr, r)

		lg.Info(
			"status", rd.status,
			"size", rd.size,
		)

	})
}
