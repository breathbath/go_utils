package http

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

var reqExample *http.Request

func init() {
	u, _ := url.Parse("http://ya.ru?someKey=someVal&someNumb=-123")
	reqExample = &http.Request{URL: u}
}

func TestGetRequestValueString(t *testing.T) {
	actualVal := GetRequestValueString(reqExample, "nonExistingParam", "someDefaultValue")
	assert.Equal(t, "someDefaultValue", actualVal)

	actualVal = GetRequestValueString(reqExample, "someKey", "someDefaultValue")
	assert.Equal(t, "someVal", actualVal)
}

func TestGetRequestValueInt(t *testing.T) {
	actualVal := GetRequestValueInt(reqExample, "nonExistingParam", 1)
	assert.EqualValues(t, 1, actualVal)

	actualVal = GetRequestValueInt(reqExample, "someNumb", 0)
	assert.EqualValues(t, -123, actualVal)
}
