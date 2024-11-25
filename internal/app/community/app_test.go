package community

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func TestGRPCServer(t *testing.T) {
	server, grpcServer, err := GetServers(&config.Config{
		DB: config.DBConnect{
			Port:    "test",
			Host:    "test",
			DBName:  "test",
			User:    "test",
			Pass:    "test",
			SSLMode: "test",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, grpcServer)
}
