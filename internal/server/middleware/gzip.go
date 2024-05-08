package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzippingResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *gzippingResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

var _ http.ResponseWriter = &LoggingResponseWriter{}

func RequestGzipper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if incoming content is encoded with gzip, decode it
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer gz.Close()
			r.Body = gz
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			grw := gzip.NewWriter(w)
			w = &gzippingResponseWriter{
				ResponseWriter: w,
				Writer:         grw,
			}
			defer grw.Close()
		}
		next.ServeHTTP(w, r)
	})
}
