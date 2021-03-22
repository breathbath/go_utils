package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildBasicAuthHeader(t *testing.T) {
	result := BuildBasicAuthString("some", "pass")
	assert.Equal(t, "c29tZTpwYXNz", result)
}
