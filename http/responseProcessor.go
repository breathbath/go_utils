package http

import "net/http"

func GetResponseStatusCode(resp *http.Response) int {
	var statusCode = 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	return statusCode
}
