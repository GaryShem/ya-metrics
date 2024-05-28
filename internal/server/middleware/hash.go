package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

var ErrHashNoMatch = errors.New("calculated and provided hashes don't match")

type HashChecker struct {
	Key string
}

func (hc *HashChecker) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hc.Key != "" && r.Header.Get("Hash") != "" {
			headerSHA := http.CanonicalHeaderKey("Hash")
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logging.Log.Errorln("Error reading body:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h := hmac.New(sha256.New, []byte(hc.Key))
			calcHashsum := base64.StdEncoding.EncodeToString(h.Sum(body))
			header := r.Header
			logging.Log.Infoln("headers: ", header)
			recHashsum := header.Get(headerSHA)
			if calcHashsum != recHashsum {
				http.Error(w, ErrHashNoMatch.Error(), http.StatusBadRequest)
				logging.Log.Errorln(ErrHashNoMatch)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r)
	})
}
