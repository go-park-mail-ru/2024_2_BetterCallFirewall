package community

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func TestGetServer(t *testing.T) {
	server, err := GetHTTPServer(&config.Config{
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
}

func TestGRPCServer(t *testing.T) {
	server, err := GetGRPCServer(&config.Config{
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
}
