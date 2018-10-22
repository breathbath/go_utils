package http

import (
	"github.com/breathbath/go_utils/env"
	"net/url"
	"strings"
)

func BuildUrl(host, path, rawQuery string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	u.Path = path
	u.RawQuery = rawQuery

	return u.String(), nil
}

func GetValidUrlFromEnvVar(urlEnvVarName string) (url.URL, error) {
	envRootUrl := env.ReadEnvOrFail(urlEnvVarName)
	parsedUrl, err := url.Parse(strings.TrimRight(envRootUrl, "/"))
	if err != nil {
		return url.URL{}, err
	}

	return *parsedUrl, nil
}
