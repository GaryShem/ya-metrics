package middleware

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
func (m *NetworkFilterMiddleware) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ips := md.Get("x-real-ip")
		if len(ips) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing x-real-ip header")
		}
		ip := net.ParseIP(ips[0])
		if ip == nil || !m.acceptedNet.Contains(ip) {
			return nil, status.Errorf(codes.Unauthenticated, "x-real-ip header value not valid")
		}
	}
	return handler(ctx, req)
}
