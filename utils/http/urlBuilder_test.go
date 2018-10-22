package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	actualUrl, err := BuildUrl("http://ya.ru", "/news", "val1=1&val2=2")
	assert.NoError(t, err)
	assert.Equal(t, "http://ya.ru/news?val1=1&val2=2", actualUrl)
}
