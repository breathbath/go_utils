package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildBasicAuthHeader(t *testing.T) {
	result := BuildBasicAuthString("some", "pass")
	assert.Equal(t, "c29tZTpwYXNz", result)
}
