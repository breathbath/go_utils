package testing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertSyncMapEqual(t *testing.T, expected, actual *sync.Map, msgAndArgs ...interface{}) bool {
	mapExpected, mapActual := make(map[interface{}]interface{}), make(map[interface{}]interface{})

	expected.Range(func(key, value interface{}) bool {
		mapExpected[key] = value
		return true
	})

	actual.Range(func(key, value interface{}) bool {
		mapActual[key] = value
		return true
	})

	return assert.Equal(t, mapExpected, mapActual, msgAndArgs...)
}
