package middleware

import (
	"hiv_mind/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rd := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rd,
		}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logger.Get().Info("http request completed",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", rd.status),
			zap.Int("size", rd.size),
			zap.Duration("duration", duration),
		)
	})
}
