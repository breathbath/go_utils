package http

import (
	"net/url"
	"strings"

	"github.com/breathbath/go_utils/v2/env"
)

func BuildURL(host, path, rawQuery string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	u.Path = path
	u.RawQuery = rawQuery

	return u.String(), nil
}

func JoinURL(parts ...string) string {
	preparedParts := make([]string, 0, len(parts))
	for k, part := range parts {
		if k == 0 {
			preparedParts = append(preparedParts, strings.TrimRight(part, "/"))
		} else {
			preparedParts = append(preparedParts, strings.Trim(part, "/"))
		}
	}
	return strings.Join(preparedParts, "/")
}

func GetValidURLFromEnvVar(urlEnvVarName string) (url.URL, error) {
	envRootURL, err := env.ReadEnvOrError(urlEnvVarName)
	if err != nil {
		return url.URL{}, err
	}

	parsedURL, err := url.Parse(strings.TrimRight(envRootURL, "/"))
	if err != nil {
		return url.URL{}, err
	}

	return *parsedURL, nil
}
