package http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUrl(t *testing.T) {
	actualURL, err := BuildURL("http://ya.ru", "/news", "val1=1&val2=2")
	assert.NoError(t, err)
	assert.Equal(t, "http://ya.ru/news?val1=1&val2=2", actualURL)
}

func TestBuildWrongUrl(t *testing.T) {
	actualURL, err := BuildURL(":slsl:", "", "")
	assert.Error(t, err)
	assert.Equal(t, "", actualURL)
}

func TestGetValidUrlFromEnvVar(t *testing.T) {
	err := os.Setenv("SOME_URL", "localhost:8080/lala?mama=1")
	assert.NoError(t, err)

	err = os.Setenv("SOME_BAD_URL", ":lsls")
	assert.NoError(t, err)

	actualURL, err := GetValidURLFromEnvVar("SOME_URL")
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080/lala?mama=1", actualURL.String())

	_, err = GetValidURLFromEnvVar("SOME_BAD_URL")
	assert.Error(t, err)

	_, err = GetValidURLFromEnvVar("SOME_NON_EXISTING_URL")
	assert.Error(t, err)

	err = os.Unsetenv("SOME_URL")
	assert.NoError(t, err)

	err = os.Unsetenv("SOME_BAD_URL")
	assert.NoError(t, err)
}

func TestJoinUrl(t *testing.T) {
	testCases := [][]string{
		{
			"http://ya.ru",
			"lala",
			"http://ya.ru/lala",
		},
		{
			"//ya.ru/",
			"dada",
			"//ya.ru/dada",
		},
		{
			"mama",
			"mama",
		},
		{
			"/papa////",
			"/papa",
		},
		{
			"https://ya.ru////",
			"//one//",
			"two",
			"/three/",
			"four/",
			"/five",
			"https://ya.ru/one/two/three/four/five",
		},
	}

	for _, testCase := range testCases {
		if len(testCase) < 2 {
			t.Errorf("Wrong strings count %d in test cases, expected amount is > 2", len(testCase))
			return
		}

		expectedResult := testCase[len(testCase)-1]
		parts := testCase[0 : len(testCase)-1]
		actualResult := JoinURL(parts...)

		assert.Equal(t, expectedResult, actualResult)
	}
}
