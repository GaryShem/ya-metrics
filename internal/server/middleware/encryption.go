package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"os"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

func NewEncryptionMiddleware(keyPath string) (*EncryptionMiddleware, error) {
	logging.Log.Infoln("initializing decryption middleware")
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &EncryptionMiddleware{
		Key: key,
	}, nil
}

// EncryptionMiddleware - middleware to support asymmetric encryption
type EncryptionMiddleware struct {
	Key []byte
}

func (hc *EncryptionMiddleware) Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logging.Log.Errorln("Error reading body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		privateKeyBlock, _ := pem.Decode(hc.Key)
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			logging.Log.Errorln("Error parsing private key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error parsing private key"))
			return
		}
		decodedBody, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, body)
		if err != nil {
			logging.Log.Errorln("Error decrypting body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error decrypting body"))
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(decodedBody))
		next.ServeHTTP(w, r)
	})
}
