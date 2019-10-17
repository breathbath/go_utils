package testing

import (
	"strings"
	"testing"
)

//AssertRequestEqual asserts that the expectations match the actual error
func AssertErrorContains(t *testing.T, theError error, msg string) {
	if theError == nil {
		t.Errorf("An error containing '%s' is expected but nil is received", msg)
		return
	}

	if strings.Contains(theError.Error(), msg) {
		return
	}

	t.Errorf("An error containing '%s' is expected but error '%s' is received", msg, theError.Error())
}
