package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttp(t *testing.T) {
	m, err := NewHTTPMetrics("http")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	m.IncErrors("test", "Ok", "")
	m.IncHits("test", "Ok", "")
	m.ObserveTiming("test", "Ok", "", 34.12)
}
