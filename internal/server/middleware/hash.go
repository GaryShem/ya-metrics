package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

var ErrHashNoMatch = errors.New("calculated and provided hashes don't match")

type HashChecker struct {
	Key string
}

func (hc *HashChecker) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hc.Key != "" {
			headerSHA := http.CanonicalHeaderKey("HashSHA256")
			body := make([]byte, r.ContentLength)
			_, err := r.Body.Read(body)
			if err != nil {
				logging.Log.Errorln("Error reading body:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h := hmac.New(sha256.New, []byte(hc.Key))
			calcHashsum := h.Sum(body)
			recHashsum := []byte(r.Header[headerSHA][0])
			if len(calcHashsum) != len(recHashsum) {
				http.Error(w, ErrHashNoMatch.Error(), http.StatusBadRequest)
				logging.Log.Errorln(ErrHashNoMatch)
			}
			for i := 0; i < len(calcHashsum); i++ {
				if calcHashsum[i] != recHashsum[i] {
					http.Error(w, ErrHashNoMatch.Error(), http.StatusBadRequest)
					logging.Log.Errorln(ErrHashNoMatch)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
