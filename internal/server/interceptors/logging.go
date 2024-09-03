package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

// ErrorLoggingInterceptor - interceptor to log grpc errors
type ErrorLoggingInterceptor struct {
}

func (m *ErrorLoggingInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	result, err := handler(ctx, req)
	if err != nil {
		logging.Log.Warnln(err)
	}
	return result, err
}
