package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

func TestGetHttpServer(t *testing.T) {
	server, err := GetHTTPServer(&config.Config{}, &metrics.HttpMetrics{})
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestGetGrpcServer(t *testing.T) {
	server, err := GetGRPCServer(&config.Config{}, &metrics.GrpcMetrics{})
	assert.NoError(t, err)
	assert.NotNil(t, server)
}
