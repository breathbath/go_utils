package url

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlBuilding(t *testing.T) {
	testCases := []struct{
		partsToGive []string
		expectedUrl string
	} {
		{
			partsToGive: []string{
				"one",
				"two",
			},
			expectedUrl: "one/two",
		},
		{
			partsToGive: []string{
				"http://ya.ru///",
				"/two",
			},
			expectedUrl: "http://ya.ru/two",
		},
		{
			partsToGive: []string{
				"/one/",
				"/two/",
			},
			expectedUrl: "/one/two",
		},
		{
			partsToGive: []string{
				"/o//ne/",
				"/two/",
			},
			expectedUrl: "/o//ne/two",
		},
		{
			partsToGive: []string{
				"some",
			},
			expectedUrl: "some",
		},
		{
			partsToGive: []string{
				"some",
				"",
			},
			expectedUrl: "some",
		},
		{
			partsToGive: []string{
				"",
			},
			expectedUrl: "",
		},
	}

	for _, testCase := range testCases {
		actualUrl := JoinURL(testCase.partsToGive...)
		assert.Equal(t, testCase.expectedUrl, actualUrl)
	}
}

