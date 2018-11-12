package rest

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

type BadResponseCodeError struct {
	url string
	resp *http.Response
}

func NewBadResponseCodeError(url string, resp *http.Response) BadResponseCodeError {
	return BadResponseCodeError{url, resp}
}

func (brce BadResponseCodeError) Error() string {
	body, readErr := ioutil.ReadAll(brce.resp.Body)
	if readErr != nil {
		body = []byte{}
	}

	return fmt.Sprintf(
		"Remote server under %s responded with code: %d, body: %s, status: %s",
		brce.url,
		brce.resp.StatusCode,
		string(body),
		brce.resp.Status,
	)
}