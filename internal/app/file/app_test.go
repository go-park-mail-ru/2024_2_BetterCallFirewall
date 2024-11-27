package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

func TestGetServer(t *testing.T) {
	server, err := GetServer(&config.Config{
		DB: config.DBConnect{
			Port:    "test",
			Host:    "test",
			DBName:  "test",
			User:    "test",
			Pass:    "test",
			SSLMode: "test",
		},
	}, &metrics.FileMetrics{})
	assert.NoError(t, err)
	assert.NotNil(t, server)
}
