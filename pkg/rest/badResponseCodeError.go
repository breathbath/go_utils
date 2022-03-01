package rest

import (
	"fmt"
	"net/http"
)

type BadResponseCodeError struct {
	url      string
	resp     *http.Response
	respBody []byte
}

func NewBadResponseCodeError(url string, resp *http.Response, respBody []byte) BadResponseCodeError {
	return BadResponseCodeError{url, resp, respBody}
}

func (brce BadResponseCodeError) Error() string {
	return fmt.Sprintf(
		"Remote server under %s responded with code: %d, body: %q, status: %s",
		brce.url,
		brce.resp.StatusCode,
		string(brce.respBody),
		brce.resp.Status,
	)
}
