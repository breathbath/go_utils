package math

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRound(t *testing.T) {
	testSets := [][]float64{
		{3.14, 0, 3},
		{3.67, 0, 4},
		{3.6788, 2, 3.68},
	}

	for _, testSet := range testSets {
		result := Round(testSet[0], testSet[1])
		assert.EqualValues(t, testSet[2], result)
	}
}

func TestCountDecimalPlaces(t *testing.T) {
	testSets := [][]float64{
		{3.14, 2},
		{3.6788, 4},
		{3, 0},
		{1.0, 0},
	}

	for _, testSet := range testSets {
		result := CountDecimalPlaces(testSet[0])
		assert.EqualValues(t, testSet[1], result)
	}
}

func TestRandInt(t *testing.T) {
	output := RandInt(10)
	assert.True(t, output < 11, "The rand int number %d should be less than %d", output, 11)
}

func TestCountDigits(t *testing.T) {
	testSets := [][]int64{
		{22, 2},
		{444444, 6},
		{0, 1},
		{-10, 2},
	}

	for _, testSet := range testSets {
		result := CountDigits(testSet[0])
		assert.EqualValues(
			t,
			testSet[1],
			result,
			"Number %d has %d digits rather than %d",
			testSet[0],
			result,
			testSet[1],
		)
	}
}
