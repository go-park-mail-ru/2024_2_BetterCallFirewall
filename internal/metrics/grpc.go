package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type GrpcMetrics struct {
	HitsTotal *prometheus.CounterVec
	name      string
	up        bool
	Timings   *prometheus.HistogramVec
	Errors    *prometheus.CounterVec
}

func NewGrpcMetrics(name string) (*GrpcMetrics, error) {
	var metric GrpcMetrics
	metric.name = name
	metric.HitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_hits_total",
			Help: "Number of total hits",
		},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metric.HitsTotal); err != nil {
		return nil, err
	}

	metric.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_total_timings",
			Buckets: []float64{0, 0.1, 0.5, 1, 5},
		},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metric.Timings); err != nil {
		return nil, err
	}

	metric.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "Number of total errors",
		},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metric.Errors); err != nil {
		return nil, err
	}
	metric.up = true
	return &metric, nil
}

func (m *GrpcMetrics) IncErrors(path string) {
	m.Errors.WithLabelValues(path, m.name).Inc()
}

func (m *GrpcMetrics) IncHits(path string) {
	m.HitsTotal.WithLabelValues(path, m.name).Inc()
}

func (m *GrpcMetrics) ObserveTiming(path string, time float64) {
	m.Timings.WithLabelValues(path, m.name).Observe(time)
}
