package middleware

import (
	"net"
	"net/http"
)

// NetworkFilterMiddleware - middleware to ascertain whether the request is coming from a trusted network
type NetworkFilterMiddleware struct {
	acceptedNet net.IPNet
}

func NewNetworkFilterMiddleware(netCIDR string) (*NetworkFilterMiddleware, error) {
	_, subnet, err := net.ParseCIDR(netCIDR)
	if err != nil {
		return nil, err
	}
	return &NetworkFilterMiddleware{acceptedNet: *subnet}, nil
}

func (m *NetworkFilterMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIPStr := r.Header.Get("X-Real-IP")
		requestIP := net.ParseIP(requestIPStr)
		if requestIP == nil || !m.acceptedNet.Contains(requestIP) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
