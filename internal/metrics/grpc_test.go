package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGrpc(t *testing.T) {
	m, err := NewGrpcMetrics("grpc")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	m.IncErrors("test")
	m.ObserveTiming("test", 13.2)
	m.IncHits("test")
}
