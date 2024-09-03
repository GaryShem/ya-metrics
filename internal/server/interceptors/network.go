package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

// NetworkFilterInterceptor - middleware to ascertain whether the request is coming from a trusted network
type NetworkFilterInterceptor struct {
	acceptedNet net.IPNet
}

func NewNetworkFilterInterceptor(netCIDR string) (*NetworkFilterInterceptor, error) {
	_, subnet, err := net.ParseCIDR(netCIDR)
	if err != nil {
		return nil, err
	}
	return &NetworkFilterInterceptor{acceptedNet: *subnet}, nil
}

func (m *NetworkFilterInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ips := md.Get("x-real-ip")
		if len(ips) == 0 {
			logging.Log.Infoln("request with no ip data")
			return nil, status.Errorf(codes.Unauthenticated, "missing x-real-ip header")
		}
		ip := net.ParseIP(ips[0])
		if ip == nil || !m.acceptedNet.Contains(ip) {
			logging.Log.Infoln("request from untrusted network")
			return nil, status.Errorf(codes.Unauthenticated, "x-real-ip header value not valid")
		}
	}
	return handler(ctx, req)
}
