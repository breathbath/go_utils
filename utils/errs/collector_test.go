package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	actualError := CollectErrors(",", err1, err2)
	assert.EqualError(t, actualError, "error 1,error 2")

	assert.Nil(t, CollectErrors(",", nil))
	assert.Nil(t, CollectErrors(","))
}

func TestWrapError(t *testing.T) {
	err := errors.New("some internal error")
	outputErr := WrapError(err, ",", "some external error %d", 1)
	assert.EqualError(t, outputErr, "some external error 1,some internal error")

	assert.EqualError(t, WrapError(nil, ",", "some external error %d", 1), "some external error 1")
}

func TestAppendError(t *testing.T) {
	errs := []error{}
	AppendError(errors.New("some error"), &errs)
	assert.Len(t, errs, 1)

	errs = []error{}
	AppendError(nil, &errs)
	assert.Len(t, errs, 0)
}

func TestErrorContainer(t *testing.T) {
	errCont := NewErrorContainer()
	assert.Nil(t, errCont.Result(","))

	errCont.AddError(errors.New("new error"))
	errCont.AddErrorF("Some formatted error %s, %d", "err_text", 1)
	assert.EqualError(t, errCont.Result(","), "new error,Some formatted error err_text, 1")
}
