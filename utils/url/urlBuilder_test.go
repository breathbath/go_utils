package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlBuilding(t *testing.T) {
	testCases := []struct {
		partsToGive []string
		expectedURL string
	}{
		{
			partsToGive: []string{
				"one",
				"two",
			},
			expectedURL: "one/two",
		},
		{
			partsToGive: []string{
				"http://ya.ru///",
				"/two",
			},
			expectedURL: "http://ya.ru/two",
		},
		{
			partsToGive: []string{
				"/one/",
				"/two/",
			},
			expectedURL: "/one/two",
		},
		{
			partsToGive: []string{
				"/o//ne/",
				"/two/",
			},
			expectedURL: "/o//ne/two",
		},
		{
			partsToGive: []string{
				"some",
			},
			expectedURL: "some",
		},
		{
			partsToGive: []string{
				"some",
				"",
			},
			expectedURL: "some",
		},
		{
			partsToGive: []string{
				"",
			},
			expectedURL: "",
		},
	}

	for _, testCase := range testCases {
		actualURL := JoinURL(testCase.partsToGive...)
		assert.Equal(t, testCase.expectedURL, actualURL)
	}
}
