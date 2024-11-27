package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	size int64
	res  string
}

func TestGetSizeRange(t *testing.T) {
	tests := []TestCase{
		{size: 511, res: "0 - 0.5 MB"},
		{size: 1000 * 1024, res: "0.5 - 1 MB"},
		{size: 1024 * 1024 * 9, res: "1 - 10 MB"},
		{size: 1024 * 1024 * 11, res: "> 10 MB"},
	}

	for _, test := range tests {
		actual := getSizeRange(test.size)
		if actual != test.res {
			t.Errorf("getSizeRange(%d) = %s, want %s", test.size, actual, test.res)
		}
	}
}

func TestNewMetrics(t *testing.T) {
	m, err := NewFileMetrics("file")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	m.IncHits("test", "200", "GET", "format", 10)
	m.ObserveTiming("test", "200", "GET", "format", 10, 12)
	m.IncErrors("test", "200", "GET", "format", 10)
}
