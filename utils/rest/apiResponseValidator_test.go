package rest

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSuccessResponses(t *testing.T) {
	successResponseCodes := []int{200, 201, 301, 304, 205}
	for _, respCode := range successResponseCodes {
		resp := http.Response{StatusCode: respCode}
		assert.Nil(t, ValidateResponse("someurl", &resp, []byte{}))
	}
}

func TestFailureResponses(t *testing.T) {
	failedResponseCodes := []int{500, 400, 403, 404}
	for _, respCode := range failedResponseCodes {
		resp := http.Response{StatusCode: respCode, Status: fmt.Sprint(respCode), Body: ioutil.NopCloser(bytes.NewBufferString("Hello World")),}
		actualError := ValidateResponse("someurl", &resp, []byte("Hello World"))
		expectedErrText := fmt.Sprintf("Remote server under someurl responded with code: %d, body: %s, status: %d", respCode, "Hello World", respCode)
		assert.EqualError(t, actualError, expectedErrText)
	}
}
