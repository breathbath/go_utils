package http

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func getReqExample(t *testing.T) *http.Request {
	correctURL, err := url.Parse("http://ya.ru?someKey=someVal&someNumb=-123&time=2001-01-01T11-00-00")
	require.NoError(t, err)

	return &http.Request{URL: correctURL}
}

func getReqExampleWrong(t *testing.T) *http.Request {
	wrongURL, err := url.Parse("http://ya.ru?time=2001-01-01")
	require.NoError(t, err)
	return &http.Request{URL: wrongURL}
}

func TestGetRequestValueString(t *testing.T) {
	actualVal := GetRequestValueString(getReqExample(t), "nonExistingParam", "someDefaultValue")
	assert.Equal(t, "someDefaultValue", actualVal)

	actualVal = GetRequestValueString(getReqExample(t), "someKey", "someDefaultValue")
	assert.Equal(t, "someVal", actualVal)
}

func TestGetRequestValueInt(t *testing.T) {
	actualVal := GetRequestValueInt(getReqExample(t), "nonExistingParam", 1)
	assert.EqualValues(t, 1, actualVal)

	actualVal = GetRequestValueInt(getReqExample(t), "someNumb", 0)
	assert.EqualValues(t, -123, actualVal)

	actualVal = GetRequestValueInt(getReqExample(t), "someKey", -1)
	assert.EqualValues(t, -1, actualVal)
}

func TestGetRequestValueTimeWithError(t *testing.T) {
	actualTime, err := GetRequestValueTimeWithError(getReqExample(t), "time")
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	assert.Equal(t, expectedTime.UTC(), actualTime)

	_, err = GetRequestValueTimeWithError(getReqExampleWrong(t), "time")
	assert.EqualError(t, err, `parsing time "2001-01-01" as "2006-01-02T15-04-05": cannot parse "" as "T"`)

	_, err = GetRequestValueTimeWithError(getReqExample(t), "lala")
	assert.EqualError(t, err, `no time value provided for key lala`)
}

func TestGetRequestValueTimeWithDefaultValue(t *testing.T) {
	defaultValue, err := time.Parse("2006-01-02T15:04:05", "2002-02-02T12:02:02")
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	actualTime := GetRequestValueTime(getReqExample(t), "time", defaultValue)
	assert.Equal(t, expectedTime, actualTime)

	actualTime = GetRequestValueTime(getReqExampleWrong(t), "time", defaultValue)
	assert.Equal(t, defaultValue, actualTime)
}
