package io

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type mockCloser struct {
	wasClosed bool
	errToGive error
}

func (mc *mockCloser) Close() error {
	mc.wasClosed = true

	return mc.errToGive
}

func TestClosingNilResource(t *testing.T) {
	var nilCloser io.Closer
	CloseResourceSecure("some res", nilCloser)
}

func TestIgnoringClosingError(t *testing.T) {
	mc := &mockCloser{
		errToGive: errors.New("some err"),
	}
	CloseResourceSecure("some res", mc)
	assert.True(t, mc.wasClosed)
}

func TestSuccessfulClose(t *testing.T) {
	mc := &mockCloser{
		errToGive: nil,
	}
	CloseResourceSecure("some res", mc)
	assert.True(t, mc.wasClosed)
}
