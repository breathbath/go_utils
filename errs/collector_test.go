package errs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectErrors(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")

	actualError := CollectErrors(",", err1, err2)
	assert.EqualError(t, actualError, "Error 1,Error 2")

	assert.Nil(t, CollectErrors(",", nil))
	assert.Nil(t, CollectErrors(","))
}

func TestWrapError(t *testing.T) {
	err := errors.New("Some internal error")
	outputErr := WrapError(err, ",", "Some external error %d", 1)
	assert.EqualError(t, outputErr, "Some external error 1,Some internal error")

	assert.EqualError(t, WrapError(nil, ",", "Some external error %d", 1), "Some external error 1")
}

func TestAppendError(t *testing.T) {
	errs := []error{}
	AppendError(errors.New("Some error"), &errs)
	assert.Len(t, errs, 1)

	errs = []error{}
	AppendError(nil, &errs)
	assert.Len(t, errs, 0)
}

func TestErrorContainer(t *testing.T) {
	errCont := NewErrorContainer()
	assert.Nil(t, errCont.Result(","))

	errCont.AddError(errors.New("New error"))
	errCont.AddErrorF("Some formatted error %s, %d", "err_text", 1)
	assert.EqualError(t, errCont.Result(","), "New error,Some formatted error err_text, 1")
}
