package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponseStatusCode(t *testing.T) {
	resp := http.Response{StatusCode: 200}
	assert.Equal(t, 200, GetResponseStatusCode(&resp))

	assert.Equal(t, 0, GetResponseStatusCode(nil))
}
