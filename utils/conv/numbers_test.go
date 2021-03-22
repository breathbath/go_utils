package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input          string
	defaultVal     int64
	expectedOutput int64
}

func TestExtractIntFromString(t *testing.T) {
	testCases := []testCase{
		{
			"123",
			0,
			123,
		},
		{
			"nana233",
			0,
			233,
		},
		{
			"naoa",
			10,
			10,
		},
	}

	for _, tc := range testCases {
		actualOutput := ExtractIntFromString(tc.input, tc.defaultVal)
		assert.Equal(t, tc.expectedOutput, actualOutput)
	}
}

func TestConvertFloatToLongStringNumber(t *testing.T) {
	var flVal float64 = 2.123456789123456
	actualOutput := ConvertFloatToLongStringNumber(flVal)
	assert.Equal(t, "2.123456789123456", actualOutput)
}
