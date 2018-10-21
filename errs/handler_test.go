package errs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailOnErrorWithExistingError(t *testing.T) {
	assert.PanicsWithValue(t, "Some err", func() {
		err := errors.New("Some err")
		FailOnError(err)
	})
}

func TestFailOnErrorWithEmptyError(t *testing.T) {
	FailOnError(nil)
}

func TestFailOnErrorF(t *testing.T) {
	assert.PanicsWithValue(t, "Some error: fail", func() {
		FailOnErrorF("Some error: %s", "fail")
	})
}
