package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

type RequestData struct {
	uri      string
	method   string
	execTime time.Duration

	writer *LoggingResponseWriter
}

type (
	responseData struct {
		status int
		size   int
		body   string
	}
	LoggingResponseWriter struct {
		http.ResponseWriter
		data *responseData
	}
)

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.data.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *LoggingResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.data.size += size
	w.data.body = string(data)
	return size, err
}

var _ http.ResponseWriter = &LoggingResponseWriter{}

func NewRequestData(w http.ResponseWriter, r *http.Request) *RequestData {
	return &RequestData{
		uri:      r.RequestURI,
		method:   r.Method,
		execTime: 0,
		writer: &LoggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: http.StatusOK,
				size:   0,
				body:   "",
			},
		},
	}
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		data := NewRequestData(w, r)
		next.ServeHTTP(data.writer, r)
		data.execTime = time.Since(startTime)
		logging.Log.Infoln(
			"uri", data.uri,
			"method", data.method,
			"statusCode", data.writer.data.status,
			"size", data.writer.data.size,
			"resLength", data.writer.data.size,
			"execTime", data.execTime,
			fmt.Sprintf("headers %#v", w.Header()),
			"body", data.writer.data.body,
		)
	})
}
