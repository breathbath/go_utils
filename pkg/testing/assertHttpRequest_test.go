package testing

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestEquality(t *testing.T) {
	req := &http.Request{
		Method:     http.MethodPut,
		URL:        &url.URL{Path: "/some"},
		Proto:      "",
		ProtoMajor: 0,
		ProtoMinor: 0,
		Header: http.Header{
			"Authorisation": []string{"Bearer 123"},
			"Content Type":  []string{"application/json"},
		},
		Body:             io.NopCloser(bytes.NewBufferString("Some body")),
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}

	localT := &testing.T{}
	AssertRequestEqual(localT, &RequestExpectation{IsNilRequestExpected: true}, nil)
	assert.False(t, localT.Failed())

	localT = &testing.T{}
	AssertRequestEqual(localT, &RequestExpectation{IsNilRequestExpected: true}, req)
	assert.True(t, localT.Failed())

	localT = &testing.T{}
	AssertRequestEqual(localT, &RequestExpectation{IsNilRequestExpected: false}, nil)
	assert.True(t, localT.Failed())

	localT = &testing.T{}
	AssertRequestEqual(
		localT,
		&RequestExpectation{
			IsNilRequestExpected: false,
			ExpectedURL:          NewStringExpectation("/some"),
			ExpectedMethod:       NewStringExpectation(http.MethodPut),
			ExpectedHeaders: []KeyValueStr{
				{
					Key:   "Authorisation",
					Value: "Bearer 123",
				},
				{
					Key:   "Content Type",
					Value: "application/json",
				},
			},
			ExpectedBody: NewStringExpectation("Some body"),
		},
		req,
	)
	assert.False(t, localT.Failed())
}
