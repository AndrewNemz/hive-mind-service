package middleware

import (
	compressdata "hiv_mind/pkg/compress_data"
	"hiv_mind/pkg/logger"
	"net/http"
	"strings"
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

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := compressdata.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := compressdata.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			r.ContentLength = -1
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
