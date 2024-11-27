package ext_grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetGRPCProvider(t *testing.T) {
	cl, err := GetGRPCProvider("", "")
	assert.NoError(t, err)
	assert.NotNil(t, cl)
}
