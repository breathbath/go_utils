package http

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	actualUrl, err := BuildUrl("http://ya.ru", "/news", "val1=1&val2=2")
	assert.NoError(t, err)
	assert.Equal(t, "http://ya.ru/news?val1=1&val2=2", actualUrl)
}

func TestBuildWrongUrl(t *testing.T) {
	actualUrl, err := BuildUrl(":slsl:", "", "")
	assert.Error(t, err)
	assert.Equal(t, "", actualUrl)
}

func TestGetValidUrlFromEnvVar(t *testing.T) {
	err := os.Setenv("SOME_URL", "localhost:8080/lala?mama=1")
	assert.NoError(t, err)

	err = os.Setenv("SOME_BAD_URL", ":lsls")
	assert.NoError(t, err)

	actualUrl, err := GetValidUrlFromEnvVar("SOME_URL")
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080/lala?mama=1", actualUrl.String())

	_, err = GetValidUrlFromEnvVar("SOME_BAD_URL")
	assert.Error(t, err)

	_, err = GetValidUrlFromEnvVar("SOME_NON_EXISTING_URL")
	assert.Error(t, err)

	err = os.Unsetenv("SOME_URL")
	assert.NoError(t, err)

	err = os.Unsetenv("SOME_BAD_URL")
	assert.NoError(t, err)
}