package http

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var reqExample *http.Request
var reqExampleWrong *http.Request

func init() {
	correctUrl, _ := url.Parse("http://ya.ru?someKey=someVal&someNumb=-123&time=2001-01-01T11-00-00")
	wrongUrl, _ := url.Parse("http://ya.ru?time=2001-01-01")
	reqExample = &http.Request{URL: correctUrl}
	reqExampleWrong = &http.Request{URL: wrongUrl}
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

func TestGetRequestValueTimeWithError(t *testing.T) {
	actualTime , err := GetRequestValueTimeWithError(reqExample, "time")
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	assert.Equal(t, expectedTime.UTC(), actualTime)

	_, err = GetRequestValueTimeWithError(reqExampleWrong, "time")
	assert.EqualError(t, err, `parsing time "2001-01-01" as "2006-01-02T15-04-05": cannot parse "" as "T"`)
}

func TestGetRequestValueTimeWithDefaultValue(t *testing.T) {
	defaultValue, err := time.Parse("2006-01-02T15:04:05", "2002-02-02T12:02:02")
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	actualTime := GetRequestValueTime(reqExample, "time", defaultValue)
	assert.Equal(t, expectedTime, actualTime)

	actualTime = GetRequestValueTime(reqExampleWrong, "time", defaultValue)
	assert.Equal(t, defaultValue, actualTime)
}
