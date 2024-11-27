package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig("")
	assert.NotNil(t, err)
	assert.Nil(t, cfg)

	cfg, err = GetConfig("./test.env")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
