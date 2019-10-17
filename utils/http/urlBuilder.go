package http

import (
	"github.com/breathbath/go_utils/utils/env"
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

func JoinUrl(parts ...string) string {
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

func GetValidUrlFromEnvVar(urlEnvVarName string) (url.URL, error) {
	envRootUrl, err := env.ReadEnvOrError(urlEnvVarName)
	if err != nil {
		return url.URL{}, err
	}

	parsedUrl, err := url.Parse(strings.TrimRight(envRootUrl, "/"))
	if err != nil {
		return url.URL{}, err
	}

	return *parsedUrl, nil
}
