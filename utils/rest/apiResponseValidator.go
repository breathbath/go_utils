package rest

import (
	"net/http"
)

func ValidateResponse(url string, resp *http.Response) error {
	if resp.StatusCode > 199 && resp.StatusCode < 400 {
		return nil
	}

	return NewBadResponseCodeError(url, resp)
}
