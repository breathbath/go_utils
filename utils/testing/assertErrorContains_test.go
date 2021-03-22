package testing

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertErrorContains(t *testing.T) {
	actualError := errors.New("some big error")

	localT := &testing.T{}

	AssertErrorContains(localT, actualError, "big")
	assert.False(t, localT.Failed())

	AssertErrorContains(localT, actualError, "small")
	assert.True(t, localT.Failed())

	AssertErrorContains(localT, nil, "small")
	assert.True(t, localT.Failed())
}
