package testing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapEquality(t *testing.T) {
	mapLeft := &sync.Map{}
	mapLeft.Store("red", 1)
	mapLeft.Store("blue", 2)

	mapRight := &sync.Map{}
	mapRight.Store("red", 1)
	mapRight.Store("blue", 2)

	localT := &testing.T{}
	result1 := AssertSyncMapEqual(localT, &sync.Map{}, &sync.Map{})
	assert.True(t, result1)

	result2 := AssertSyncMapEqual(localT, mapLeft, mapRight)
	assert.True(t, result2)

	mapRight.Store("yellow", 3)
	result3 := AssertSyncMapEqual(localT, mapLeft, mapRight)
	assert.False(t, result3)
}
