package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

type GrpcMiddleware struct {
	metrics *metrics.GrpcMetrics
}

func NewGrpcMiddleware(metrics *metrics.GrpcMetrics) *GrpcMiddleware {
	return &GrpcMiddleware{
		metrics: metrics,
	}
}

func (m *GrpcMiddleware) GrpcMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := next(ctx, req)
	if err != nil {
		m.metrics.IncErrors(info.FullMethod)
	}
	m.metrics.IncHits(info.FullMethod)
	m.metrics.ObserveTiming(info.FullMethod, time.Since(start).Seconds())
	return h, err
}
