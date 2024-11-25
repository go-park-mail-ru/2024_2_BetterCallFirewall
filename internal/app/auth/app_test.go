package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func TestGetHttpServer(t *testing.T) {
	server, err := GetHTTPServer(&config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestGetGrpcServer(t *testing.T) {
	server := GetGRPCServer(&config.Config{})
	assert.NotNil(t, server)
}
