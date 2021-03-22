package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailOnErrorWithExistingError(t *testing.T) {
	assert.PanicsWithValue(t, "some err", func() {
		err := errors.New("some err")
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
