package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HttpMetrics struct {
	Errors      *prometheus.CounterVec
	serviceName string
	Hits        *prometheus.CounterVec
	Timings     *prometheus.HistogramVec
}

func NewHTTPMetrics(serviceName string) (*HttpMetrics, error) {
	var metrics HttpMetrics
	metrics.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Number of total errors.",
		},
		[]string{"path", "service", "status", "method"},
	)
	if err := prometheus.Register(metrics.Errors); err != nil {
		return nil, err
	}

	metrics.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits_total",
			Help: "Number of total hits.",
		},
		[]string{"path", "service", "status", "method"},
	)
	if err := prometheus.Register(metrics.Hits); err != nil {
		return nil, err
	}

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "total_timings",
			Buckets: []float64{0, 0.1, 0.5, 1, 5},
		},
		[]string{"path", "status", "method"},
	)
	if err := prometheus.Register(metrics.Timings); err != nil {
		return nil, err
	}

	metrics.serviceName = serviceName

	return &metrics, nil
}

func (m *HttpMetrics) IncErrors(path string, status, method string) {
	m.Errors.WithLabelValues(path, m.serviceName, status, method).Inc()
}

func (m *HttpMetrics) IncHits(path string, status, method string) {
	m.Hits.WithLabelValues(path, m.serviceName, status, method).Inc()
}

func (m *HttpMetrics) ObserveTiming(path string, status, method string, time float64) {
	m.Timings.WithLabelValues(path, status, method).Observe(time)
}
