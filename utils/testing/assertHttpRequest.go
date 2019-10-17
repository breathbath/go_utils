package testing

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type KeyValueStr struct {
	Key   string
	Value string
}

//RequestExpectation contains expectation info for a http request
type RequestExpectation struct {
	IsNilRequestExpected bool
	ExpectedURL          StringExpectation
	ExpectedMethod       StringExpectation
	ExpectedHeaders      []KeyValueStr
	ExpectedBody         StringExpectation
}

//StringExpectation defines string assertion rule, if IsActive is false expectation is disabled
type StringExpectation struct {
	Value    string
	IsActive bool
}

//NewStringExpectation creates active expectation
func NewStringExpectation(value string) StringExpectation {
	return StringExpectation{Value: value, IsActive: true}
}

//AssertRequestEqual asserts that the expectations match the actual request
func AssertRequestEqual(t *testing.T, re RequestExpectation, actualRequest *http.Request) {
	if re.IsNilRequestExpected && actualRequest == nil {
		return
	}

	if actualRequest == nil {
		assert.Fail(t, "Actual request is nil, but no nil request is expected")
		return
	}

	if re.IsNilRequestExpected {
		assert.Fail(t, "Actual request is not nil, but a nil request is expected")
		return
	}

	if re.ExpectedURL.IsActive {
		assert.Equal(t, re.ExpectedURL.Value, actualRequest.URL.String())
	}

	if re.ExpectedMethod.IsActive {
		assert.Equal(t, re.ExpectedMethod.Value, actualRequest.Method)
	}

	if len(re.ExpectedHeaders) > 0 {
		for _, expHeader := range re.ExpectedHeaders {
			actualHeaderVal := actualRequest.Header.Get(expHeader.Key)
			assert.Equal(t, expHeader.Value, actualHeaderVal)
		}
	}

	if re.ExpectedBody.IsActive {
		bodyBytes, err := ioutil.ReadAll(actualRequest.Body)
		assert.NoError(t, err)
		if err == nil {
			actualRequest.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			assert.Equal(t, re.ExpectedBody.Value, string(bodyBytes))
		}
	}
}
