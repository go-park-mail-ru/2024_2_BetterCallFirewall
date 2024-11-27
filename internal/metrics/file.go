package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type FileMetrics struct {
	Errors      *prometheus.CounterVec
	serviceName string
	up          bool
	Hits        *prometheus.CounterVec
	Timings     *prometheus.HistogramVec
}

func NewFileMetrics(serviceName string) (*FileMetrics, error) {
	var metrics FileMetrics
	metrics.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_errors_total",
			Help: "Number of total errors.",
		},
		[]string{"path", "service", "status", "method", "format", "size"},
	)
	if err := prometheus.Register(metrics.Errors); err != nil {
		return nil, err
	}

	metrics.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_hits_total",
			Help: "Number of total hits.",
		},
		[]string{"path", "service", "status", "method", "format", "size"},
	)
	if err := prometheus.Register(metrics.Hits); err != nil {
		return nil, err
	}

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "file_total_timings",
			Buckets: []float64{0, 0.1, 0.5, 1, 5},
		},
		[]string{"path", "status", "method", "format", "weight"},
	)
	if err := prometheus.Register(metrics.Timings); err != nil {
		return nil, err
	}

	metrics.serviceName = serviceName
	metrics.up = true

	return &metrics, nil
}

func (m *FileMetrics) IncErrors(path string, status, method, format string, size int64) {
	newPath := pathConverter(path)
	m.Errors.WithLabelValues(newPath, m.serviceName, status, method, format, getSizeRange(size)).Inc()
}

func (m *FileMetrics) IncHits(path string, status, method, format string, size int64) {
	newPath := pathConverter(path)
	m.Hits.WithLabelValues(newPath, m.serviceName, status, method, format, getSizeRange(size)).Inc()
}

func (m *FileMetrics) ObserveTiming(path string, status, method, format string, size int64, time float64) {
	newPath := pathConverter(path)
	m.Timings.WithLabelValues(newPath, status, method, format, getSizeRange(size)).Observe(time)
}

func (m *FileMetrics) ShutDown() {
	m.up = false
}

func getSizeRange(size int64) string {
	switch {
	case size <= 511*1024:
		return "0 - 0.5 MB"
	case size <= 1023*1024:
		return "0.5 - 1 MB"
	case size <= 10*1024*1024:
		return "1 - 10 MB"
	default:
		return "> 10 MB"
	}
}
